// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"embed"
	"fmt"
	"io/fs"
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render/tmpl"
)

//go:embed testdata/*
var testFS embed.FS

func TestEmbedLoader_Load(t *testing.T) {
	t.Run("InvalidRootDirectory", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      "nonexistent",
			Extension: ".html",
		}
		loader := NewEmbed(testFS, config)
		templates, err := loader.Load("")

		assert.NotNil(t, err)
		assertTemplates(t, templates, []string{})
	})

	t.Run("EmptyExtensionMatchesAllFiles", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root: "testdata/templates",
		}
		expected := []string{
			"about.html",
			"index.html",
			"ignore.txt",
			"layouts/base.html",
			"users/profile.html",
		}

		loader := NewEmbed(testFS, config)
		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, expected)
	})

	t.Run("BasicHtmlTemplates", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      "testdata/templates",
			Extension: ".html",
		}
		expected := []string{
			"index.html",
			"about.html",
			"layouts/base.html",
			"users/profile.html",
		}

		loader := NewEmbed(testFS, config)
		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, expected)
	})
}

func TestEmbedLoader_Read(t *testing.T) {
	t.Run("ReadOutsideRootDirectory", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root: "testdata/templates",
		}

		loader := NewEmbed(testFS, config)
		content, err := loader.Read("../outside.html")

		assert.NotNil(t, err)
		assert.Len(t, content, 0)
	})

	t.Run("ReadNonexistentTemplate", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root: "testdata/templates",
		}

		loader := NewEmbed(testFS, config)
		content, err := loader.Read("nonexistent.html")

		assert.NotNil(t, err)
		assert.Len(t, content, 0)
	})

	t.Run("ReadExistingTemplate", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root: "testdata/templates",
		}

		loader := NewEmbed(testFS, config)
		content, err := loader.Read("index.html")

		assert.Nil(t, err)
		assert.Len(t, content, 83)
	})

	t.Run("ReadEmbeddedTemplateError", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root: "testdata/templates",
		}
		errFs := &errorEmbedLoader{
			readErr: fmt.Errorf("failed to read"),
		}
		loader := NewEmbed(errFs, config)
		_, err := loader.Read("test.html")

		assert.NotNil(t, err)
	})
}

func TestEmbedLoader_Extension(t *testing.T) {
	t.Run("WithExtension", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      t.TempDir(),
			Extension: ".html",
		}

		loader := NewEmbed(testFS, config)
		assert.Equals(t, loader.Extension(), ".html")
	})

	t.Run("EmptyExtension", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      t.TempDir(),
			Extension: "",
		}

		loader := NewEmbed(testFS, config)
		assert.Equals(t, loader.Extension(), "")
	})
}

type errorEmbedLoader struct {
	readErr error
	openErr error
	fs      embed.FS
}

func (l *errorEmbedLoader) Open(name string) (fs.File, error) {
	f, _ := l.fs.Open(name)
	if l.openErr != nil {
		return f, l.openErr
	}
	return f, nil
}

func (l *errorEmbedLoader) ReadFile(string) ([]byte, error) {
	if l.readErr != nil {
		return nil, l.readErr
	}
	return []byte{}, nil
}
