// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"testing"

	"github.com/nanoninja/assert"
)

func TestNew(t *testing.T) {
	tpl := New()

	assert.NotNil(t, tpl)
	assert.NotNil(t, tpl.funcMap)
	assert.Nil(t, tpl.loader)
}
