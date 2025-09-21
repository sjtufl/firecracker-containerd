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
	VMIDAnnotationKey = "aws.firecracker.vm.id"

	// VMVCpuCountAnnotationKey is the key specified in an OCI-runtime config annotation section
	// specifying the number of vCPUs to allocate to the VM.
	VMVCpuCountAnnotationKey = "aws.firecracker.vm.vcpu-count"

	// VMMemSizeMibAnnotationKey is the key specified in an OCI-runtime config annotation section
	// specifying the amount of memory (in MiB) to allocate to the VM.
	VMMemSizeMibAnnotationKey = "aws.firecracker.vm.mem-size-mib"
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
