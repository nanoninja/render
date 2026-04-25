// Copyright 2026 The Nanoninja Authors. All rights reserved.
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

// countedErrorWriter fails after n successful writes.
type countedErrorWriter struct {
	remaining int
}

func (w *countedErrorWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errors.New("write error")
	}
	w.remaining--
	return len(p), nil
}

func TestContentTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		renderer Renderer
		expected string
	}{
		{"JSON", JSON(), ContentTypeJSON},
		{"XML", XML(), ContentTypeXML},
		{"Text", Text(), ContentTypeText},
		{"HTML", HTML(), ContentTypeHTML},
		{"CSV", CSV(), ContentTypeCSV},
		{"Binary", Binary(), ContentTypeBinary},
		{"YAML", YAML(), ContentTypeYAML},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.renderer.ContentType())
		})
	}
}

func TestCheckContext(t *testing.T) {
	t.Run("ReturnsNilForActiveContext", func(t *testing.T) {
		ctx := context.Background()
		err := CheckContext(ctx)

		assert.NoError(t, err)
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
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestWithHelpers(t *testing.T) {
	t.Run("WithTimeoutSetsTimeout", func(t *testing.T) {
		opts := WithTimeout(NoOptions, 5*time.Second)
		assert.Equal(t, 5*time.Second, opts.Timeout)
	})

	t.Run("WithPrettyEnablesPretty", func(t *testing.T) {
		opts := WithPretty(NoOptions)
		assert.True(t, opts.Pretty)
	})

	t.Run("WithIndentSetsIndent", func(t *testing.T) {
		opts := WithIndent(NoOptions, "  ")
		assert.Equal(t, "  ", opts.Indent)
	})

	t.Run("WithHeaderAddsHeader", func(t *testing.T) {
		opts := WithHeader(NoOptions, "X-Test", "value")
		assert.Equal(t, "value", opts.Headers["X-Test"][0])
	})

	t.Run("WithHeaderAllocatesMap", func(t *testing.T) {
		opts := WithHeader(Options{}, "X-Test", "value")
		assert.NotNil(t, opts.Headers)
	})

	t.Run("HelpersDoNotMutateOriginal", func(t *testing.T) {
		original := NoOptions
		_ = WithTimeout(original, 5*time.Second)
		assert.Equal(t, time.Duration(0), original.Timeout)
	})
}

func TestApplyTimeout(t *testing.T) {
	t.Run("ReturnsOriginalContextWhenNoTimeout", func(t *testing.T) {
		ctx := context.Background()
		derived, cancel := ApplyTimeout(ctx, NoOptions)
		defer cancel()

		assert.Equal(t, ctx, derived)
	})

	t.Run("ReturnsDeadlineContextWhenTimeoutSet", func(t *testing.T) {
		ctx := context.Background()
		opts := Options{Timeout: 5 * time.Second}
		derived, cancel := ApplyTimeout(ctx, opts)
		defer cancel()

		_, ok := derived.Deadline()
		assert.True(t, ok)
	})
}
