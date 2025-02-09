// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"fmt"
	"testing"

	"github.com/nanoninja/assert"
)

// assertTemplates verifies that the loaded templates match the expected ones.
// It performs an order-independent comparison.
func assertTemplates(t *testing.T, got, want []string) {
	t.Helper()

	// Verify the number of templates
	assert.Equals(t, len(got), len(want), "template count mismatch")

	// Create map for efficient lookup
	paths := make(map[string]bool)
	for _, path := range got {
		paths[path] = true
	}
	// Check for expected templates
	for _, expectedPath := range want {
		assert.True(t,
			paths[expectedPath],
			fmt.Sprintf("missing expected template %s", expectedPath),
		)
	}
	// Check for unexpected templates
	for _, actualPath := range got {
		found := false
		for _, expectedPath := range want {
			if actualPath == expectedPath {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("found unexpected template: %s", actualPath))
	}
}

type errorLoader struct {
	loadErr error
	readErr error
}

func (l *errorLoader) Load(string) ([]string, error) {
	if l.loadErr != nil {
		return nil, l.loadErr
	}
	return []string{"test.html"}, nil
}

func (l *errorLoader) Read(string) ([]byte, error) {
	if l.readErr != nil {
		return nil, l.readErr
	}
	return []byte("content"), nil
}

func (l *errorLoader) Extension() string {
	return ".html"
}
