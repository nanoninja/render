// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"fmt"
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render/tmpl"
)

func TestCompositeLoader_Load(t *testing.T) {
	defaultTemplates := map[string]string{
		"index.html":   "<h1>Default Index</h1>",
		"about.html":   "<h1>Default About</h1>",
		"contact.html": "<h1>Default Contact</h1>",
		"ignored.txt":  "Ignored file",
	}
	customTemplates := map[string]string{
		"index.html":   "<h1>Custom Index</h1>",
		"about.html":   "<h1>Custom About</h1>",
		"special.html": "<h1>Special Page</h1>",
	}

	t.Run("CompositeWithPriority", func(t *testing.T) {
		config := tmpl.LoaderConfig{Extension: ".html"}
		loaders := []tmpl.Loader{
			NewString(customTemplates, config),
			NewString(defaultTemplates, config),
		}
		expected := []string{
			"index.html",
			"about.html",
			"special.html",
			"contact.html",
		}
		loader := NewComposite(loaders, config)
		templates, err := loader.Load("")

		assert.Nil(t, err)
		assert.Equals(t, loader.Extension(), config.Extension)
		assertTemplates(t, templates, expected)
	})

	t.Run("EmptyLoadersList", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Extension: ".html",
		}
		loader := NewComposite([]tmpl.Loader{}, config)
		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, []string{})
	})

	t.Run("DifferentExtensions", func(t *testing.T) {
		loaders := []tmpl.Loader{
			NewString(map[string]string{
				"test.gohtml": "Custom template",
			}, tmpl.LoaderConfig{Extension: ".gohtml"}),
			NewString(map[string]string{
				"test.html": "Default template",
			}, tmpl.LoaderConfig{Extension: ".html"}),
		}
		expected := []string{
			"test.gohtml",
			"test.html",
		}
		config := tmpl.LoaderConfig{Extension: ""}
		loader := NewComposite(loaders, config)

		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, expected)
	})

	t.Run("LoaderWithError", func(t *testing.T) {
		loaders := []tmpl.Loader{
			&errorLoader{loadErr: fmt.Errorf("load failed")},
		}
		config := tmpl.LoaderConfig{}
		loader := NewComposite(loaders, config)

		_, err := loader.Load("")
		assert.NotNil(t, err)
	})
}

func TestCompositeLoader_Read(t *testing.T) {
	defaultTemplates := map[string]string{
		"index.html":   "<h1>Default Index</h1>",
		"about.html":   "<h1>Default About</h1>",
		"contact.html": "<h1>Default Contact</h1>",
		"ignored.txt":  "Ignored file",
	}
	customTemplates := map[string]string{
		"index.html":   "<h1>Custom Index</h1>",
		"about.html":   "<h1>Custom About</h1>",
		"special.html": "<h1>Special Page</h1>",
	}

	t.Run("CompositeWithPriority", func(t *testing.T) {
		config := tmpl.LoaderConfig{Extension: ".html"}
		loaders := []tmpl.Loader{
			NewString(customTemplates, config),
			NewString(defaultTemplates, config),
		}
		expected := map[string]string{
			"index.html":   "<h1>Custom Index</h1>",
			"contact.html": "<h1>Default Contact</h1>",
			"special.html": "<h1>Special Page</h1>",
		}
		loader := NewComposite(loaders, config)

		_, err := loader.Load("")
		assert.Nil(t, err)

		for expecteName, expectedContent := range expected {
			content, err := loader.Read(expecteName)

			assert.Nil(t, err)
			assert.Equals(t, string(content), expectedContent)
		}
	})

	t.Run("ReadError", func(t *testing.T) {
		loaders := []tmpl.Loader{
			&errorLoader{readErr: fmt.Errorf("read failed")},
		}
		config := tmpl.LoaderConfig{}
		loader := NewComposite(loaders, config)

		_, err := loader.Read("test.html")
		assert.NotNil(t, err)
	})

	t.Run("TemplateNotFoundError", func(t *testing.T) {
		loaders := []tmpl.Loader{
			&errorLoader{readErr: tmpl.ErrTemplateNotFound},
		}
		config := tmpl.LoaderConfig{}
		loader := NewComposite(loaders, config)

		_, err := loader.Read("test.html")
		assert.ErrorIs(t, err, tmpl.ErrTemplateNotFound)
	})
}
