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

package cloner

import (
	"context"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// New returns a new cloner.
func New(depth int, stdout io.Writer) Cloner {
	return &cloner{
		depth:  depth,
		stdout: stdout,
	}
}

// NewDefault returns a cloner with default settings.
func NewDefault() Cloner {
	return New(1, os.Stdout)
}

// default cloner using the built-in Git client.
type cloner struct {
	depth  int
	stdout io.Writer
}

// Clone the repository using the built-in Git client.
func (c *cloner) Clone(ctx context.Context, params Params) error {
	opts := &git.CloneOptions{
		RemoteName: "origin",
		Progress:   c.stdout,
		URL:        params.Repo,
	}
	// set the reference name if provided
	if params.Ref != "" {
		opts.ReferenceName = plumbing.ReferenceName(expandRef(params.Ref))
	}
	// set depth if cloning the head commit of a branch as
	// opposed to a specific commit sha
	if params.Sha == "" {
		opts.Depth = c.depth
	}

	// clone the repository
	r, err := git.PlainClone(params.Dir, false, opts)
	if err != nil {
		return err
	}
	if params.Sha == "" {
		return nil
	}

	// checkout the sha
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(params.Sha),
	})
}
