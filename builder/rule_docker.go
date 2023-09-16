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

// ConfigureDocker configures a Docker step.
func ConfigureDocker(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check if a Dockerfile exists
	// TODO check subdirectories with glob
	if !exists(fsys, "Dockerfile") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	useImage := isContainerRuntime(pipeline)

	// add the docker build step
	{
		repo := "hello/world" // dummy name
		// TODO parse the .git/config and get the remote orign
		// url. extract the repository name from the url and use
		// this as the image name, if possible.

		script := new(spec.StepPlugin)
		script.Image = "plugins/docker"
		script.With = map[string]interface{}{
			"tags":    "latest",
			"repo":    repo,
			"dry_run": true,
		}

		if useImage {
			script.Image = "plugins/docker"
			script.Privileged = true
		} else {
			// TODO we should eventually use the container-less
			// version of the plugin here
			script.Image = "plugins/docker"
		}

		step := new(spec.Step)
		step.Name = "docker_build"
		step.Type = "plugin"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
