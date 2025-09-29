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

package firecrackeroci

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
)

const (
	// VMIDAnnotationKey is the key specified in an OCI-runtime config annotation section
	// specifying the ID of the VM in which the container should be spun up.
	VMIDAnnotationKey = "firecracker.vm.id"

	// VMVCpuCountAnnotationKey is the key specified in an OCI-runtime config annotation section
	// specifying the number of vCPUs to allocate to the VM.
	VMVCpuCountAnnotationKey = "firecracker.vm.vcpu-count"

	// VMMemSizeMibAnnotationKey is the key specified in an OCI-runtime config annotation section
	// specifying the amount of memory (in MiB) to allocate to the VM.
	VMMemSizeMibAnnotationKey = "firecracker.vm.mem-size-mib"

	// Additional clone related annotations
	VMProvisionMode      = "firecracker.vm.provision-mode"
	VMSnapshotMemoryPath = "firecracker.vm.snapshot.memory-path"
	VMSnapshotStatePath  = "firecracker.vm.snapshot.vmstate-path"

	// JSON map of device ID to path on host. Example:
	// {"MN2HE43UOVRDA": "/dev/vdb", "NQ2HE43UOVRDA": "/dev/vdc"}
	// This is used to override the drives specified in the snapshot.
	// If this annotation is not provided, the drives from the snapshot
	// will be used as-is.
	VMSnapshotDriveOverrides = "firecracker.vm.snapshot.drive-overrides"
)

type ProvisionMode string

const (
	// ProvisionModeClone indicates that the VM should be created by cloning from a snapshot.
	ProvisionModeClone ProvisionMode = "Clone"
	// ProvisionModeCreate indicates that the VM should be created from scratch.
	ProvisionModeCreate ProvisionMode = "Create"
)

// WithVMID annotates a containerd client's container object with a given firecracker VMID.
func WithVMID(vmID string) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		if s.Annotations == nil {
			s.Annotations = make(map[string]string)
		}

		s.Annotations[VMIDAnnotationKey] = vmID
		return nil
	}
}

func WithVMSizeConfig(vcpuCount, memSizeMib uint32) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		if s.Annotations == nil {
			s.Annotations = make(map[string]string)
		}

		if vcpuCount > 0 {
			s.Annotations[VMVCpuCountAnnotationKey] = fmt.Sprintf("%d", vcpuCount)
		}

		if memSizeMib > 0 {
			s.Annotations[VMMemSizeMibAnnotationKey] = fmt.Sprintf("%d", memSizeMib)
		}

		return nil
	}
}

func WithVMProvisionMode(mode ProvisionMode, snapshotMemPath, snapshotStatePath string) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		if s.Annotations == nil {
			s.Annotations = make(map[string]string)
		}

		s.Annotations[VMProvisionMode] = string(mode)
		if mode == ProvisionModeClone {
			if snapshotMemPath == "" || snapshotStatePath == "" {
				return fmt.Errorf("snapshot memory path and state path must be provided when using clone provision mode")
			}
			s.Annotations[VMSnapshotMemoryPath] = snapshotMemPath
			s.Annotations[VMSnapshotStatePath] = snapshotStatePath
		}

		return nil
	}
}

// ValidateVMProvisionMode checks that the VM provision mode annotation is valid.
func ValidateVMProvisionMode(mode string) error {
	if mode != string(ProvisionModeCreate) && mode != string(ProvisionModeClone) && mode != "" {
		return fmt.Errorf("invalid VM provision mode %q: must be either %q or %q", mode, ProvisionModeCreate, ProvisionModeClone)
	}
	return nil
}

// WithVMSnapshotDriveOverridesJsonMapString annotates a containerd client's container object
// with a given firecracker snapshot drive overrides.
func WithVMSnapshotDriveOverridesJsonMapString(driveOverrides string) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *oci.Spec) error {
		if s.Annotations == nil {
			s.Annotations = make(map[string]string)
		}

		s.Annotations[VMSnapshotDriveOverrides] = driveOverrides
		return nil
	}
}
