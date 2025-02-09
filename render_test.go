// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

type errorWriterTest struct{}

func (*errorWriterTest) Write([]byte) (int, error) {
	return 0, errors.New("error writer test")
}

func TestMimesTypes(t *testing.T) {
	tests := []struct {
		name     string
		opt      func(*Options)
		expected string
	}{
		{
			name:     "MimeTextPlain",
			opt:      MimeTextPlain(),
			expected: "text/plain; charset=utf-8",
		},
		{
			name:     "MimeTextHTML",
			opt:      MimeTextHTML(),
			expected: "text/html; charset=utf-8",
		},
		{
			name:     "MimeJSON",
			opt:      MimeJSON(),
			expected: "application/json; charset=utf-8",
		},
		{
			name:     "MimeXML",
			opt:      MimeXML(),
			expected: "application/xml; charset=utf-8",
		},
		{
			name:     "MimeCSV",
			opt:      MimeCSV(),
			expected: "text/csv; charset=utf-8",
		},
		{
			name:     "MimeBinary",
			opt:      MimeBinary(),
			expected: "application/octet-stream",
		},
		{
			name:     "MimePDF",
			opt:      MimePDF(),
			expected: "application/pdf",
		},
		{
			name:     "MimeStream",
			opt:      MimeStream(),
			expected: "application/octet-stream",
		},
		{
			name:     "MimeYAML",
			opt:      MimeYAML(),
			expected: "application/yaml; charset=utf-8",
		},
	}

	for _, tt := range tests {
		opts := NewOptions()
		opts.Use(tt.opt)

		assert.Equals(t, opts.ContentType(), tt.expected)
	}
}

func TestCheckContext(t *testing.T) {
	t.Run("ReturnsNilForActiveContext", func(t *testing.T) {
		ctx := context.Background()
		err := CheckContext(ctx)

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorForCancelledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := CheckContext(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("ReturnsErrorForTimeoutContext", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond)

		err := CheckContext(ctx)
		assert.Equals(t, err, context.DeadlineExceeded)
	})
}
