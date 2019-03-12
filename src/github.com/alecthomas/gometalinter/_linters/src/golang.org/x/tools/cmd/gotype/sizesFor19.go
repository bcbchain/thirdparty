// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

package main

import "go/types"

const defaultCompiler = "source"

func SizesFor(compiler, arch string) types.Sizes {
	return types.SizesFor(compiler, arch)
}
