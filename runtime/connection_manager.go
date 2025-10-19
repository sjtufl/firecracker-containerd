// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	taskAPI "github.com/containerd/containerd/api/runtime/task/v2"
	"github.com/containerd/ttrpc"
	"github.com/firecracker-microvm/firecracker-go-sdk/vsock"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/firecracker-microvm/firecracker-containerd/eventbridge"
	agentTaskTtrpc "github.com/firecracker-microvm/firecracker-containerd/proto/service/agenttask/ttrpc"
	drivemount "github.com/firecracker-microvm/firecracker-containerd/proto/service/drivemount/ttrpc"
	ioproxy "github.com/firecracker-microvm/firecracker-containerd/proto/service/ioproxy/ttrpc"
)

// ConnectionConfig holds configuration for vsock connection management
type ConnectionConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
	ConnectTimeout time.Duration
}

// DefaultConnectionConfig returns default connection configuration
func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
		BackoffFactor:  2.0,
		ConnectTimeout: defaultVSockConnectTimeout,
	}
}

// AgentClients holds all the agent client interfaces
type AgentClients struct {
	TaskClient        taskAPI.TaskService
	AgentTaskClient   agentTaskTtrpc.AgentTaskService
	EventBridgeClient eventbridge.Getter
	DriveMountClient  drivemount.DriveMounterService
	IOProxyClient     ioproxy.IOProxyService
}

// ConnectionManager manages the vsock connection and provides resilient client access
type ConnectionManager struct {
	vsockPath string
	vsockPort uint32
	logger    *logrus.Entry
	config    *ConnectionConfig

	mu         sync.RWMutex
	conn       net.Conn
	rpcClient  *ttrpc.Client
	clients    *AgentClients
	connecting bool
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(vsockPath string, vsockPort uint32, logger *logrus.Entry, config *ConnectionConfig) *ConnectionManager {
	if config == nil {
		config = DefaultConnectionConfig()
	}

	return &ConnectionManager{
		vsockPath: vsockPath,
		vsockPort: vsockPort,
		logger:    logger,
		config:    config,
	}
}

// Connect establishes the initial vsock connection and creates clients
func (cm *ConnectionManager) Connect(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.clients != nil {
		return nil // Already connected
	}

	return cm.connect(ctx)
}

// connect performs the actual connection (must be called with lock held)
func (cm *ConnectionManager) connect(ctx context.Context) error {
	cm.logger.WithFields(logrus.Fields{
		"vsock_path": cm.vsockPath,
		"vsock_port": cm.vsockPort,
	}).Debug("establishing vsock connection")

	conn, err := vsock.DialContext(ctx, cm.vsockPath, cm.vsockPort, vsock.WithLogger(cm.logger))
	if err != nil {
		return fmt.Errorf("failed to dial VM over vsock: %w", err)
	}

	rpcClient := ttrpc.NewClient(conn, ttrpc.WithOnClose(func() { _ = conn.Close() }))

	clients := &AgentClients{
		TaskClient:        taskAPI.NewTaskClient(rpcClient),
		AgentTaskClient:   agentTaskTtrpc.NewAgentTaskClient(rpcClient),
		EventBridgeClient: eventbridge.NewGetterClient(rpcClient),
		DriveMountClient:  drivemount.NewDriveMounterClient(rpcClient),
		IOProxyClient:     ioproxy.NewIOProxyClient(rpcClient),
	}

	cm.conn = conn
	cm.rpcClient = rpcClient
	cm.clients = clients

	cm.logger.Debug("vsock connection established successfully")
	return nil
}

// GetClients returns the current clients, attempting reconnection if needed
func (cm *ConnectionManager) GetClients(ctx context.Context) (*AgentClients, error) {
	cm.mu.RLock()
	if cm.clients != nil {
		clients := cm.clients
		cm.mu.RUnlock()
		return clients, nil
	}
	cm.mu.RUnlock()

	// Need to connect or reconnect
	return cm.reconnectIfNeeded(ctx)
}

// reconnectIfNeeded attempts to reconnect with retry logic
func (cm *ConnectionManager) reconnectIfNeeded(ctx context.Context) (*AgentClients, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Double-check pattern - another goroutine might have connected
	if cm.clients != nil {
		return cm.clients, nil
	}

	// Avoid multiple concurrent connection attempts
	if cm.connecting {
		cm.mu.Unlock()
		// Wait a bit and retry
		time.Sleep(50 * time.Millisecond)
		cm.mu.Lock()
		if cm.clients != nil {
			return cm.clients, nil
		}
		return nil, fmt.Errorf("connection attempt in progress")
	}

	cm.connecting = true
	defer func() { cm.connecting = false }()

	// Attempt reconnection with exponential backoff
	backoff := cm.config.InitialBackoff
	for attempt := 0; attempt <= cm.config.MaxRetries; attempt++ {
		if attempt > 0 {
			cm.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"backoff": backoff,
			}).Warn("retrying vsock connection")

			// Sleep with context cancellation support
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}

			// Increase backoff for next attempt
			nextBackoff := time.Duration(float64(backoff) * cm.config.BackoffFactor)
			if nextBackoff > cm.config.MaxBackoff {
				nextBackoff = cm.config.MaxBackoff
			}
			backoff = nextBackoff
		}

		connectCtx, cancel := context.WithTimeout(ctx, cm.config.ConnectTimeout)
		err := cm.connect(connectCtx)
		cancel()

		if err == nil {
			return cm.clients, nil
		}

		cm.logger.WithError(err).WithField("attempt", attempt+1).Warn("vsock connection failed")

		// Clean up failed connection attempt
		cm.cleanup()
	}

	return nil, fmt.Errorf("failed to establish vsock connection after %d attempts", cm.config.MaxRetries+1)
}

// MarkConnectionBroken marks the current connection as broken and cleans it up
func (cm *ConnectionManager) MarkConnectionBroken() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.clients != nil {
		cm.logger.Warn("marking vsock connection as broken")
		cm.cleanup()
	}
}

// cleanup closes the connection and clears clients (must be called with lock held)
func (cm *ConnectionManager) cleanup() {
	if cm.rpcClient != nil {
		cm.rpcClient.Close()
		cm.rpcClient = nil
	}
	if cm.conn != nil {
		cm.conn.Close()
		cm.conn = nil
	}
	cm.clients = nil
}

// Close closes the connection manager
func (cm *ConnectionManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cleanup()
	return nil
}

// isConnectionError determines if an error indicates a broken connection
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific error types that indicate connection issues
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}

	// Check for gRPC status codes that indicate connection issues
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled:
			return true
		}
	}

	// Check for specific error messages that indicate connection problems
	errStr := err.Error()
	connectionErrorPatterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"EOF",
		"connection lost",
		"no such file or directory", // vsock path not found
		"context deadline exceeded: unknown",
	}

	for _, pattern := range connectionErrorPatterns {
		if len(errStr) > 0 && len(pattern) > 0 &&
			len(errStr) >= len(pattern) && errStr[len(errStr)-len(pattern):] == pattern {
			return true
		}
	}

	return false
}
