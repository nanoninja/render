// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loader

import (
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render/tmpl"
)

func TestStringLoader_Load(t *testing.T) {
	templates := map[string]string{
		"index.html":     "<h1>Index</h1>",
		"about.html":     "<h1>About</h1>",
		"contact.txt":    "Contact info",
		"profile.gohtml": "Profile template",
	}

	t.Run("BasicHtmlTemplates", func(t *testing.T) {
		expected := []string{
			"about.html",
			"index.html",
		}
		config := tmpl.LoaderConfig{
			Extension: ".html",
		}

		loader := NewString(templates, config)
		names, err := loader.Load("")

		assert.NoError(t, err)
		assert.Equal(t, config.Extension, loader.Extension())
		assertTemplates(t, names, expected)
	})

	t.Run("EmptyExtension", func(t *testing.T) {
		expected := []string{
			"about.html",
			"contact.txt",
			"index.html",
			"profile.gohtml",
		}
		config := tmpl.LoaderConfig{
			Extension: "",
		}

		loader := NewString(templates, config)
		names, err := loader.Load("")

		assert.NoError(t, err)
		assert.Equal(t, config.Extension, loader.Extension())
		assertTemplates(t, names, expected)
	})

	t.Run("EmptyTemplates", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Extension: ".html",
		}

		loader := NewString(map[string]string{}, config)
		names, err := loader.Load("")

		assert.NoError(t, err)
		assert.Equal(t, config.Extension, loader.Extension())
		assertTemplates(t, names, []string{})
	})
}

func TestStringLoader_Read(t *testing.T) {
	templates := map[string]string{
		"index.html":     "<h1>Index</h1>",
		"about.html":     "<h1>About</h1>",
		"contact.txt":    "Contact info",
		"profile.gohtml": "Profile template",
	}

	t.Run("ReadTemplates", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Extension: ".html",
		}
		expected := map[string]string{
			"index.html": "<h1>Index</h1>",
			"about.html": "<h1>About</h1>",
		}
		loader := NewString(templates, config)

		_, err := loader.Load("")
		assert.NoError(t, err)

		for name, html := range expected {
			content, err := loader.Read(name)

			assert.NoError(t, err)
			assert.Equal(t, html, string(content))
		}
	})

	t.Run("MissingTemplate", func(t *testing.T) {
		config := tmpl.LoaderConfig{}
		loader := NewString(templates, config)

		_, err := loader.Load("")
		assert.NoError(t, err)

		_, err = loader.Read("nonexistent.html")
		assert.Error(t, err)
		assert.ErrorIs(t, err, tmpl.ErrTemplateNotFound)
	})
}
