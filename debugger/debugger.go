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

// Package debugger debugs a Harness pipelines
package debugger

import (
	"io/fs"
)

type (
	Input struct {
		// Config provides the pipeline configuration
		// encoded as a yaml document.
		Config []byte

		// Logs provides the pipeline logs which can be
		// used to identify the problem.
		Logs []byte

		// Error provides the pipeline error which can
		// be used to identify the problem.
		Error []byte

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
		// Diagnosis provides a diagnosis of the problem
		// in text or markdown format.
		Diagnosis []byte

		// Patch provides a diff / patch that can be applied
		// to the codebase to fix the problem.
		Patch []byte
	}
)

// Debugger debugs a yaml file.
type Debugger interface {
	Debug(in *Input) (*Output, error)
}
