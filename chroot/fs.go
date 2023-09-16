// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package chroot implements a chroot virtual file system.
package chroot

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ensure io/fs interface conformance conformance.
var (
	_ fs.FS         = (*FS)(nil)
	_ fs.StatFS     = (*FS)(nil)
	_ fs.GlobFS     = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
)

// FS is a chroot file system that implements the FS
// interface in the Go standard library.
type FS struct {
	base string
}

// New returns a new choot filesystem.
func New(path string) (*FS, error) {
	base, err := filepath.Abs(path)
	return &FS{base: base}, err
}

// Open opens the named file.
func (fs *FS) Open(name string) (fs.File, error) {
	path := filepath.Join(fs.base, name)
	return os.Open(path)
}

// Stat returns a FileInfo describing the named file from
// the file system.
func (fs *FS) Stat(name string) (fs.FileInfo, error) {
	path := filepath.Join(fs.base, name)
	return os.Stat(path)
}

// Glob returns the names of all files matching pattern
// or nil if there is no matching file.
func (fs *FS) Glob(pattern string) ([]string, error) {
	pattern = filepath.Join(fs.base, pattern)
	matches, err := filepath.Glob(pattern)
	for i := 0; i < len(matches); i++ {
		// trim the base prefix from the match name.
		matches[i] = strings.TrimPrefix(matches[i], fs.base)
	}
	return matches, err
}

// ReadDir reads the named directory
func (fs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	path := filepath.Join(fs.base, name)
	return os.ReadDir(path)
}

// ReadFile reads the named file from the file system fs
// and returns its contents.
func (fs *FS) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(fs.base, name)
	return os.ReadFile(path)
}
