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

package builder

import (
	"io/fs"

	spec "github.com/bradrydzewski/spec/yaml"
)

// ConfigureDefault configures a default step if the system
// is unable to automatically add any language-specific steps.
func ConfigureDefault(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0]

	// ignore if stage already contains steps
	if len(stage.Steps) == 0 {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	var image string
	if isContainerRuntime(pipeline) {
		image = "alpine"
	}

	// add dummy hello world step
	step := createScriptStep(image, "echo", "echo hello world")
	stage.Steps = append(stage.Steps, step)

	return nil
}
