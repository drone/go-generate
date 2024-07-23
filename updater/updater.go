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

// Package updater updates a Harness pipeline
package updater

import (
	"io/fs"

	"github.com/r3labs/diff"
)

type (
	Input struct {
		// Prompt provides a user-defined prompt
		// to update the pipeline configuration.
		// e.g. add a Slack notification step.
		Prompt string

		// Config defines the pipeline configuration
		// encoded as a yaml document.
		Config []byte

		// Code provides virtual filesystem access to the
		// source code, providing the AI with context around
		// the code can assist with build and test-level
		// optimizations.
		Code fs.FS

		// we may want to include more data,
		// including connectors, secrets,
		// templates, etc so that the AI has
		// enough context to make recommendations.
	}

	Output struct {
		Before []byte
		After  []byte
		Change *diff.Change
	}
)

// Updater updates a yaml file based on the given prompt.
type Updater interface {
	Update(in *Input) (*Output, error)
}
