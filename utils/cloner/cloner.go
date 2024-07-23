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

// Package cloner provides support for cloning git repositories.
package cloner

import "context"

type (
	// Params provides clone params.
	Params struct {
		Repo string
		Ref  string
		Sha  string
		Dir  string // Target clone directory.

		// clone credentials (not yet implemented)
		Username   string
		Password   string
		Privatekey string
	}

	// Cloner clones a repository.
	Cloner interface {
		// Clone a repository.
		Clone(context.Context, Params) error
	}
)
