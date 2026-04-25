// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

// YAMLConfig defines configuration for the YAML renderer.
type YAMLConfig struct {
	// Indent specifies the number of spaces used for indentation.
	// Defaults to 2 if not set.
	Indent int
}

// yamlRenderer implements YAML rendering with configurable formatting options.
type yamlRenderer struct {
	config YAMLConfig
}

// YAML creates a new YAML renderer with safe default configuration:
// - 2-space indentation
// This is the recommended constructor for most use cases.
func YAML() Renderer {
	return NewYAML(YAMLConfig{Indent: 2})
}

// NewYAML creates a YAML renderer with custom configuration.
func NewYAML(c YAMLConfig) Renderer {
	return &yamlRenderer{config: c}
}

func (r *yamlRenderer) ContentType() string {
	return ContentTypeYAML
}

// Render writes the YAML representation of data with context support.
// It handles:
// - Custom indentation based on configuration
// - Content type setting to application/yaml
func (r *yamlRenderer) Render(ctx context.Context, w io.Writer, data any, _ Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(r.config.Indent)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("yaml: %w", err)
	}
	return nil
}
