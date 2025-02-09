// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render/tmpl"
)

func TestNewFS(t *testing.T) {
	t.Run("ValidConfiguration", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      t.TempDir(),
			Extension: ".html",
		}

		loader, err := NewFS(config)

		assert.Nil(t, err)
		assert.NotNil(t, loader)
		assert.Equals(t, loader.Extension(), config.Extension)
	})

	t.Run("DirectoryWithNoPermissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		target := filepath.Join(tmpDir, "locked")

		assert.Nil(t, os.MkdirAll(target, 0755))
		assert.Nil(t, os.Chmod(tmpDir, 0000))

		t.Cleanup(func() {
			_ = os.Chmod(tmpDir, 0755)
		})

		config := tmpl.LoaderConfig{Root: target}
		loader, err := NewFS(config)

		assert.ErrorIs(t, err, ErrInvalidRoot)
		assert.StringContains(t, err.Error(), "invalid root path: cannot access root directory")
		assert.Nil(t, loader)
	})

	t.Run("FileInsteadOfDirectory", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "file.txt")

		err := os.WriteFile(filePath, []byte("content"), 0644)
		assert.Nil(t, err)

		config := tmpl.LoaderConfig{Root: filePath}
		loader, err := NewFS(config)

		assert.ErrorIs(t, err, ErrInvalidRoot)
		assert.StringContains(t, err.Error(), "invalid root path: root path is not a directory")
		assert.Nil(t, loader)
	})

	t.Run("InvalidRootDirectory", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      "nonexistent",
			Extension: ".html",
		}

		loader, err := NewFS(config)
		assert.ErrorIs(t, err, ErrInvalidRoot)
		assert.Nil(t, loader)
	})
}

func TestLoader_Load(t *testing.T) {
	t.Run("BasicHtmlTemplates", func(t *testing.T) {
		tmpDir := t.TempDir()

		files := map[string]string{
			"index.html":            "<h1>Index</h1>",
			"about.html":            "<h1>About</h1>",
			"users/profile.html":    "<h1>Profile</h1>",
			"layouts/base.html":     "<!DOCTYPE html>",
			"ignore.txt":            "ignored file",
			"users/settings.gohtml": "wrong extension",
		}

		for path, content := range files {
			fullPath := filepath.Join(tmpDir, path)
			err := os.MkdirAll(filepath.Dir(fullPath), 0755)
			assert.Nil(t, err)

			err = os.WriteFile(fullPath, []byte(content), 0644)
			assert.Nil(t, err)
		}

		expected := []string{
			"about.html",
			"index.html",
			"layouts/base.html",
			"users/profile.html",
		}
		config := tmpl.LoaderConfig{
			Root:      tmpDir,
			Extension: ".html",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, expected)
	})

	t.Run("EmptyExtensionMatchesAllFiles", func(t *testing.T) {
		tmpDir := t.TempDir()

		files := map[string]string{
			"test.html": "<h1>Test</h1>",
			"test.txt":  "Text file",
		}

		for path, content := range files {
			err := os.WriteFile(filepath.Join(tmpDir, path), []byte(content), 0644)
			assert.Nil(t, err)
		}

		config := tmpl.LoaderConfig{
			Root:      tmpDir,
			Extension: "",
		}
		expected := []string{
			"test.html",
			"test.txt",
		}

		loader, err := NewFS(config)
		assert.Nil(t, err)

		templates, err := loader.Load("")

		assert.Nil(t, err)
		assertTemplates(t, templates, expected)
	})
}

func TestFSLoader_Read(t *testing.T) {
	t.Run("ReadOutsideRootDirectory", func(t *testing.T) {
		config := tmpl.LoaderConfig{Root: t.TempDir()}

		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Read("../index.html")
		assert.NotNil(t, err)
	})

	t.Run("ReadNonexistentFile", func(t *testing.T) {
		config := tmpl.LoaderConfig{Root: t.TempDir()}

		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Read("nonexistent.html")
		assert.NotNil(t, err)
	})

	t.Run("FileExistsButNotReadable", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmplPath := "test.html"
		fullPath := filepath.Join(tmpDir, tmplPath)

		err := os.WriteFile(fullPath, []byte("content"), 0644)
		assert.Nil(t, err)

		err = os.Chmod(fullPath, 0000)
		assert.Nil(t, err)

		t.Cleanup(func() {
			_ = os.Chmod(fullPath, 0644)
		})

		config := tmpl.LoaderConfig{Root: tmpDir}

		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Read(tmplPath)
		assert.NotNil(t, err)
	})

	t.Run("ReadExistingFile", func(t *testing.T) {
		tmpDir := t.TempDir()
		content := "template content"
		tmplPath := "test.html"

		err := os.WriteFile(
			filepath.Join(tmpDir, tmplPath),
			[]byte(content),
			0644,
		)
		assert.Nil(t, err)

		config := tmpl.LoaderConfig{Root: tmpDir}
		loader, err := NewFS(config)

		assert.Nil(t, err)

		_, err = loader.Read(tmplPath)
		assert.Nil(t, err)
	})
}

func TestFSLoader_Extension(t *testing.T) {
	t.Run("WithExtension", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      t.TempDir(),
			Extension: ".html",
		}
		loader, err := NewFS(config)

		assert.Nil(t, err)
		assert.Equals(t, loader.Extension(), ".html")
	})

	t.Run("EmptyExtension", func(t *testing.T) {
		config := tmpl.LoaderConfig{
			Root:      t.TempDir(),
			Extension: "",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		assert.Equals(t, loader.Extension(), "")
	})
}

func TestFSLoader_EdgeCases2(t *testing.T) {
	t.Run("PermissionDeniedOnSubdirectory", func(t *testing.T) {
		tmpDir := t.TempDir()
		subDir := filepath.Join(tmpDir, "subdir")

		err := os.MkdirAll(subDir, 0755)
		assert.Nil(t, err)

		err = os.WriteFile(
			filepath.Join(subDir, "test.html"),
			[]byte("test content"),
			0644,
		)

		assert.Nil(t, err)
		assert.Nil(t, os.Chmod(subDir, 0000))

		t.Cleanup(func() {
			_ = os.Chmod(subDir, 0755)
		})

		config := tmpl.LoaderConfig{
			Root:      tmpDir,
			Extension: ".html",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Load("")
		assert.NotNil(t, err)
	})

	t.Run("BrokenSymlink", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := os.Symlink(
			filepath.Join(tmpDir, "nonexistent.html"),
			filepath.Join(tmpDir, "broken.html"),
		)
		assert.Nil(t, err)

		config := tmpl.LoaderConfig{
			Root:      tmpDir,
			Extension: ".html",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Load("")
		assert.NotNil(t, err)
	})

	t.Run("SymlinkOutsideRootDirectory", func(t *testing.T) {
		tmpDir := t.TempDir()
		outsideDir := t.TempDir()

		err := os.WriteFile(
			filepath.Join(outsideDir, "external.html"),
			[]byte("external content"),
			0644,
		)
		assert.Nil(t, err)

		err = os.Symlink(
			filepath.Join(outsideDir, "external.html"),
			filepath.Join(tmpDir, "symlink.html"),
		)
		assert.Nil(t, err)

		config := tmpl.LoaderConfig{
			Root:      tmpDir,
			Extension: ".html",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Load("")
		assert.NotNil(t, err)
	})

	t.Run("InvalidRelativePath", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("This test is Windows-specific")
		}
		config := tmpl.LoaderConfig{
			Root: "C:\\templates",
		}
		loader, err := NewFS(config)
		assert.Nil(t, err)

		_, err = loader.Load("")
		assert.NotNil(t, err)
	})
}
