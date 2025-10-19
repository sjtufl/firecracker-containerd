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

	taskAPI "github.com/containerd/containerd/api/runtime/task/v2"
	apitypes "github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/protobuf/types"
	"github.com/golang/protobuf/ptypes/empty"

	agentTaskTtrpc "github.com/firecracker-microvm/firecracker-containerd/proto/service/agenttask/ttrpc"
	drivemount "github.com/firecracker-microvm/firecracker-containerd/proto/service/drivemount/ttrpc"
	ioproxy "github.com/firecracker-microvm/firecracker-containerd/proto/service/ioproxy/ttrpc"
)

// ResilientTaskClient provides a resilient wrapper around taskAPI.TaskService
type ResilientTaskClient struct {
	cm *ConnectionManager
}

// NewResilientTaskClient creates a new resilient task client
func NewResilientTaskClient(cm *ConnectionManager) *ResilientTaskClient {
	return &ResilientTaskClient{cm: cm}
}

func (r *ResilientTaskClient) State(ctx context.Context, req *taskAPI.StateRequest) (*taskAPI.StateResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.State(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Create(ctx context.Context, req *taskAPI.CreateTaskRequest) (*taskAPI.CreateTaskResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Create(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Start(ctx context.Context, req *taskAPI.StartRequest) (*taskAPI.StartResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Start(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Delete(ctx context.Context, req *taskAPI.DeleteRequest) (*taskAPI.DeleteResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Delete(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Pids(ctx context.Context, req *taskAPI.PidsRequest) (*taskAPI.PidsResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Pids(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Pause(ctx context.Context, req *taskAPI.PauseRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Pause(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Resume(ctx context.Context, req *taskAPI.ResumeRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Resume(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Checkpoint(ctx context.Context, req *taskAPI.CheckpointTaskRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Checkpoint(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Kill(ctx context.Context, req *taskAPI.KillRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Kill(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Exec(ctx context.Context, req *taskAPI.ExecProcessRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Exec(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) ResizePty(ctx context.Context, req *taskAPI.ResizePtyRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.ResizePty(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) CloseIO(ctx context.Context, req *taskAPI.CloseIORequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.CloseIO(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Update(ctx context.Context, req *taskAPI.UpdateTaskRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Update(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Wait(ctx context.Context, req *taskAPI.WaitRequest) (*taskAPI.WaitResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Wait(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Stats(ctx context.Context, req *taskAPI.StatsRequest) (*taskAPI.StatsResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Stats(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Connect(ctx context.Context, req *taskAPI.ConnectRequest) (*taskAPI.ConnectResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Connect(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientTaskClient) Shutdown(ctx context.Context, req *taskAPI.ShutdownRequest) (*types.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.TaskClient.Shutdown(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

// ResilientAgentTaskClient provides a resilient wrapper around agentTaskTtrpc.AgentTaskService
type ResilientAgentTaskClient struct {
	cm *ConnectionManager
}

// NewResilientAgentTaskClient creates a new resilient agent task client
func NewResilientAgentTaskClient(cm *ConnectionManager) *ResilientAgentTaskClient {
	return &ResilientAgentTaskClient{cm: cm}
}

func (r *ResilientAgentTaskClient) ExecuteCommand(ctx context.Context, req *agentTaskTtrpc.ExecuteCommandRequest) (*agentTaskTtrpc.ExecuteCommandResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.AgentTaskClient.ExecuteCommand(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientAgentTaskClient) ListExistingTasks(ctx context.Context, req *agentTaskTtrpc.ListExistingTasksRequest) (*agentTaskTtrpc.ListExistingTasksResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.AgentTaskClient.ListExistingTasks(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

// ResilientEventBridgeClient provides a resilient wrapper around eventbridge.Getter
type ResilientEventBridgeClient struct {
	cm *ConnectionManager
}

func NewResilientEventBridgeClient(cm *ConnectionManager) *ResilientEventBridgeClient {
	return &ResilientEventBridgeClient{cm: cm}
}

func (r *ResilientEventBridgeClient) GetEvent(ctx context.Context) (*apitypes.Envelope, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.EventBridgeClient.GetEvent(ctx)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

// ResilientDriveMountClient provides a resilient wrapper around drivemount.DriveMounterService
type ResilientDriveMountClient struct {
	cm *ConnectionManager
}

func NewResilientDriveMountClient(cm *ConnectionManager) *ResilientDriveMountClient {
	return &ResilientDriveMountClient{cm: cm}
}

func (r *ResilientDriveMountClient) MountDrive(ctx context.Context, req *drivemount.MountDriveRequest) (*empty.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.DriveMountClient.MountDrive(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientDriveMountClient) UnmountDrive(ctx context.Context, req *drivemount.UnmountDriveRequest) (*empty.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.DriveMountClient.UnmountDrive(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

// ResilientIOProxyClient provides a resilient wrapper around ioproxy.IOProxyService
type ResilientIOProxyClient struct {
	cm *ConnectionManager
}

func NewResilientIOProxyClient(cm *ConnectionManager) *ResilientIOProxyClient {
	return &ResilientIOProxyClient{cm: cm}
}

func (r *ResilientIOProxyClient) State(ctx context.Context, req *ioproxy.StateRequest) (*ioproxy.StateResponse, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.IOProxyClient.State(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}

func (r *ResilientIOProxyClient) Attach(ctx context.Context, req *ioproxy.AttachRequest) (*empty.Empty, error) {
	for attempt := 0; attempt <= 1; attempt++ {
		clients, err := r.cm.GetClients(ctx)
		if err != nil {
			return nil, err
		}

		result, err := clients.IOProxyClient.Attach(ctx, req)
		if err == nil {
			return result, nil
		}

		if isConnectionError(err) && attempt == 0 {
			r.cm.MarkConnectionBroken()
			continue
		}

		return nil, err
	}

	return nil, context.DeadlineExceeded
}
