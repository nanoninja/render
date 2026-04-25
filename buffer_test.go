// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = Buffer(nil)
	_ Renderer = NewBuffer(nil, BufferConfig{})
	_ Renderer = (*mockRenderer)(nil)
)

func TestBufferRenderer(t *testing.T) {
	t.Run("ContentType", func(t *testing.T) {
		assert.Equal(t, ContentTypeJSON, Buffer(JSON()).ContentType())
	})

	t.Run("RenderSimpleContent", func(t *testing.T) {
		var w bytes.Buffer

		renderer := &mockRenderer{content: "test content"}

		err := Buffer(renderer).Render(context.Background(), &w, nil, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "test content", w.String())
	})

	t.Run("UseInitialBufferSize", func(t *testing.T) {
		var w bytes.Buffer

		renderer := &mockRenderer{content: "test content"}
		config := BufferConfig{InitialSize: 100}

		err := NewBuffer(renderer, config).Render(context.Background(), &w, nil, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "test content", w.String())
	})

	t.Run("HandlePostProcessing", func(t *testing.T) {
		var w bytes.Buffer

		renderer := &mockRenderer{content: "test content"}
		config := BufferConfig{
			PostRender: func(content []byte) ([]byte, error) {
				return []byte(strings.ToUpper(string(content))), nil
			},
		}

		err := NewBuffer(renderer, config).Render(context.Background(), &w, nil, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "TEST CONTENT", w.String())
	})

	t.Run("HandleUderlyingRendererError", func(t *testing.T) {
		var w bytes.Buffer

		renderer := &mockRenderer{err: errors.New("render error")}

		err := Buffer(renderer).Render(context.Background(), &w, nil, NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "render error")
		assert.Equal(t, "", w.String())
	})

	t.Run("HandleWriteError", func(t *testing.T) {
		w := &errorWriterTest{}
		renderer := &mockRenderer{content: "test content"}

		err := Buffer(renderer).Render(context.Background(), w, nil, NoOptions)

		assert.Error(t, err)
	})

	t.Run("RespectContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		renderer := &mockRenderer{content: "test content"}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Buffer(renderer).Render(ctx, &w, nil, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, "", w.String())
	})

	t.Run("RespectTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Buffer(JSON()).Render(context.Background(), &w, map[string]string{"key": "value"}, Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Buffer(JSON()).Render(ctx, &w, map[string]string{"key": "value"}, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestBuffer_ContextCheckBeforePostProcess(t *testing.T) {
	var w bytes.Buffer

	renderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(content []byte) ([]byte, error) {
			return content, nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := NewBuffer(renderer, config).Render(ctx, &w, nil, NoOptions)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestBufferRenderer_PostProcessError(t *testing.T) {
	var w bytes.Buffer

	expectedErr := errors.New("post-process error")
	mockRenderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(_ []byte) ([]byte, error) {
			return nil, expectedErr
		},
	}

	err := NewBuffer(mockRenderer, config).Render(context.Background(), &w, nil, NoOptions)

	assert.ErrorIs(t, err, expectedErr)
	assert.StringContains(t, err.Error(), "post-processing")
}

func TestBufferRenderer_CancelDuringPostProcessing(t *testing.T) {
	var w bytes.Buffer
	var postProcessingCalled bool

	startPostProcess := make(chan struct{})

	renderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(content []byte) ([]byte, error) {
			postProcessingCalled = true
			close(startPostProcess)

			time.Sleep(50 * time.Millisecond)
			return content, nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-startPostProcess
		cancel()
	}()

	err := NewBuffer(renderer, config).Render(ctx, &w, nil, NoOptions)

	assert.ErrorIs(t, err, context.Canceled)
	assert.True(t, postProcessingCalled, "post-processing should have been called")
	assert.Equal(t, "", w.String())
}

func TestBufferRenderer_CancelBeforeWrite(t *testing.T) {
	var writeCalled bool
	postProcessStarted := make(chan struct{})

	mockRenderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(content []byte) ([]byte, error) {
			close(postProcessStarted)
			time.Sleep(50 * time.Millisecond)
			return content, nil
		},
	}
	w := &checkWriterTest{
		onWrite: func(p []byte) (int, error) {
			writeCalled = true
			return len(p), nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-postProcessStarted
		cancel()
	}()

	err := NewBuffer(mockRenderer, config).Render(ctx, w, nil, NoOptions)

	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, writeCalled)
}

func TestBufferRenderer_CancelBeforePostProcessing(t *testing.T) {
	var w bytes.Buffer
	var postProcessStarted bool

	mockRenderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(content []byte) ([]byte, error) {
			postProcessStarted = true
			return content, nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := NewBuffer(mockRenderer, config).Render(ctx, &w, nil, NoOptions)

	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, postProcessStarted)
	assert.Equal(t, "", w.String())
}

func TestBeforeFinalWrite(t *testing.T) {
	var writeCalled bool

	postProcessDone := make(chan struct{})

	renderer := &mockRenderer{content: "test content"}
	config := BufferConfig{
		PostRender: func(content []byte) ([]byte, error) {
			close(postProcessDone)

			time.Sleep(50 * time.Millisecond)
			return content, nil
		},
	}
	w := &checkWriterTest{
		onWrite: func(p []byte) (int, error) {
			writeCalled = true
			return len(p), nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-postProcessDone
		cancel()
	}()

	err := NewBuffer(renderer, config).Render(ctx, w, nil, NoOptions)

	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, writeCalled)
}

func TestBufferRenderer_ContextCheckDuringPostProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	renderer := &mockRenderer{
		content: "test content",
		beforeWrite: func() {
			time.Sleep(50 * time.Millisecond)
		},
	}
	config := BufferConfig{
		PostRender: func(b []byte) ([]byte, error) {
			return b, nil
		},
	}

	err := NewBuffer(renderer, config).Render(ctx, &bytes.Buffer{}, nil, NoOptions)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

type mockRenderer struct {
	content      string
	err          error
	onRenderDone func()
	beforeWrite  func()
}

func (r *mockRenderer) ContentType() string { return "" }

func (r *mockRenderer) Render(_ context.Context, w io.Writer, _ any, _ Options) error {
	if r.err != nil {
		return r.err
	}
	if r.beforeWrite != nil {
		r.beforeWrite()
	}
	_, err := w.Write([]byte(r.content))
	if r.onRenderDone != nil {
		r.onRenderDone()
	}
	return err
}

type checkWriterTest struct {
	onWrite func([]byte) (int, error)
}

func (w *checkWriterTest) Write(p []byte) (int, error) {
	return w.onWrite(p)
}
