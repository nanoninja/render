// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import "errors"

// ErrTemplateNotFound is returned when a specific template cannot be found
// in the configured root directory.
var ErrTemplateNotFound = errors.New("template not found")
