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

	spec "github.com/drone/spec/dist/go"
)

// ConfigureGo configures a Go step.
func ConfigureGo(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the go.mod file.
	if !exists(fsys, "go.mod") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	useImage := isContainerRuntime(pipeline)

	// add the go install step
	if exists(fsys, "main.go") {
		script := new(spec.StepRun)
		script.Script = []string{"go build"}

		if useImage {
			script.Container = new(spec.Container)
			script.Container.Image = "golang:1"
		}

		step := new(spec.Step)
		step.Name = "go_install"
		step.Type = "run"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	} else {
		script := new(spec.StepRun)
		script.Script = []string{"go install ./..."}

		if useImage {
			script.Container = new(spec.Container)
			script.Container.Image = "golang:1"
		}

		step := new(spec.Step)
		step.Name = "go_install"
		step.Type = "run"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the go test step
	{
		script := new(spec.StepRun)
		script.Script = []string{"go test -v ./..."}

		if useImage {
			script.Container = new(spec.Container)
			script.Container.Image = "golang:1"
		}

		step := new(spec.Step)
		step.Name = "go_test"
		step.Type = "run"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
