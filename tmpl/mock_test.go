// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"fmt"
	"io"
)

// mockLoader simulates a simple template loader for testing.
// It stores templates in memory using a map where keys are template names
// and values are template contents. This loader is used for basic test cases
// where we need to verify standard template loading and rendering behavior.
type mockLoader struct {
	templates map[string]string
}

func (m *mockLoader) Load(string) ([]string, error) {
	var names []string
	for name := range m.templates {
		names = append(names, name)
	}
	return names, nil
}

func (m *mockLoader) Read(name string) ([]byte, error) {
	if content, ok := m.templates[name]; ok {
		return []byte(content), nil
	}
	return nil, io.ErrUnexpectedEOF
}

func (m *mockLoader) Extension() string {
	return ".tmpl"
}

// stepErrorLoader provides fine-grained control over where errors occur
// during the template loading process.
// The step field determines at which point the error should occur:
// 1 = Load error,
// 2 = Read error,
// 3 = Parse error (via invalid content).
// This allows testing different error scenarios in the loading pipeline.
type stepErrorLoader struct {
	step int
}

func (l *stepErrorLoader) Load(string) ([]string, error) {
	if l.step == 1 {
		return nil, fmt.Errorf("simulated Load error")
	}
	return []string{"template.tmpl"}, nil
}

func (l *stepErrorLoader) Read(string) ([]byte, error) {
	if l.step == 2 {
		return nil, fmt.Errorf("simulated Read error")
	}
	return []byte("invalid {{ template"), nil
}

func (l *stepErrorLoader) Extension() string {
	return ".tmpl"
}
