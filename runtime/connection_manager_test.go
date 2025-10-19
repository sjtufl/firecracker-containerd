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
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestConnectionManager_Reconnection(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	config := &ConnectionConfig{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
		ConnectTimeout: 1 * time.Second,
	}

	cm := NewConnectionManager("/fake/path", 12345, logger, config)

	// Test that connection error detection works
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"connection refused", errors.New("connection refused"), true},
		{"broken pipe", errors.New("broken pipe"), true},
		{"EOF", errors.New("EOF"), true},
		{"other error", errors.New("some other error"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isConnectionError(tc.err)
			if result != tc.expected {
				t.Errorf("isConnectionError(%v) = %v, expected %v", tc.err, result, tc.expected)
			}
		})
	}

	// Test connection manager creation
	if cm == nil {
		t.Fatal("NewConnectionManager returned nil")
	}

	if cm.config != config {
		t.Error("Config not set correctly")
	}

	if cm.vsockPath != "/fake/path" {
		t.Error("VSock path not set correctly")
	}

	if cm.vsockPort != 12345 {
		t.Error("VSock port not set correctly")
	}
}

func TestConnectionConfig_Defaults(t *testing.T) {
	config := DefaultConnectionConfig()

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries=3, got %d", config.MaxRetries)
	}

	if config.InitialBackoff != 100*time.Millisecond {
		t.Errorf("Expected InitialBackoff=100ms, got %v", config.InitialBackoff)
	}

	if config.MaxBackoff != 5*time.Second {
		t.Errorf("Expected MaxBackoff=5s, got %v", config.MaxBackoff)
	}

	if config.BackoffFactor != 2.0 {
		t.Errorf("Expected BackoffFactor=2.0, got %f", config.BackoffFactor)
	}

	if config.ConnectTimeout != defaultVSockConnectTimeout {
		t.Errorf("Expected ConnectTimeout=%v, got %v", defaultVSockConnectTimeout, config.ConnectTimeout)
	}
}
