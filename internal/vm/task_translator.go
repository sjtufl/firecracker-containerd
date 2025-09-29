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

package vm

import (
	"sync"
)

// TaskTranslator provides ID translation between external (host) and internal (VM) task IDs
// This is used for clone mode where external task IDs need to map to existing tasks in the VM
type TaskTranslator interface {
	// TranslateToInternal maps external (host) task ID to internal (VM) task ID
	TranslateToInternal(externalID string) string
	// TranslateToExternal maps internal (VM) task ID to external (host) task ID
	TranslateToExternal(internalID string) string
	// RegisterCloneMapping registers a mapping for clone mode
	RegisterCloneMapping(externalID, internalID string)
	// IsCloneTask returns true if the given external ID is a clone task
	IsCloneTask(externalID string) bool
	// GetCloneMappings returns all current clone mappings (external -> internal)
	GetCloneMappings() map[string]string
}

// NewTaskTranslator creates a new task translator instance
func NewTaskTranslator() TaskTranslator {
	return &taskTranslator{
		externalToInternal: make(map[string]string),
		internalToExternal: make(map[string]string),
	}
}

type taskTranslator struct {
	mu                 sync.RWMutex
	externalToInternal map[string]string // external -> internal mapping
	internalToExternal map[string]string // internal -> external mapping
}

func (t *taskTranslator) TranslateToInternal(externalID string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if internalID, exists := t.externalToInternal[externalID]; exists {
		return internalID
	}
	// If no mapping exists, return the original ID (normal non-clone case)
	return externalID
}

func (t *taskTranslator) TranslateToExternal(internalID string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if externalID, exists := t.internalToExternal[internalID]; exists {
		return externalID
	}
	// If no mapping exists, return the original ID (normal non-clone case)
	return internalID
}

func (t *taskTranslator) RegisterCloneMapping(externalID, internalID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.externalToInternal[externalID] = internalID
	t.internalToExternal[internalID] = externalID
}

func (t *taskTranslator) IsCloneTask(externalID string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	_, exists := t.externalToInternal[externalID]
	return exists
}

func (t *taskTranslator) GetCloneMappings() map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for ext, internal := range t.externalToInternal {
		result[ext] = internal
	}
	return result
}
