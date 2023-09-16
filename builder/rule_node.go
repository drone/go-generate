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

// ConfigureNode configures a Node step.
func ConfigureNode(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the package.json file.
	if !exists(fsys, "package.json") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	var image string
	if isContainerRuntime(pipeline) {
		image = "node"
	}

	// add the npm install step
	stage.Steps = append(stage.Steps, createScriptStep(image,
		"npm_install",
		"npm install",
	))

	// parse the package.json file and unmarshal
	json := new(packageJson)
	if err := unmarshal(fsys, "package.json", &json); err != nil {
		return nil
	}

	// add well-known test
	if _, ok := json.Scripts["test"]; ok {
		stage.Steps = append(stage.Steps, createScriptStep(image,
			"npm_test",
			"npm run test",
		))
	}

	// add well-known lint command
	if _, ok := json.Scripts["lint"]; ok {
		stage.Steps = append(stage.Steps, createScriptStep(image,
			"npm_test",
			"npm run lint",
		))
	}

	// add well-known e2e command
	if _, ok := json.Scripts["e2e"]; ok {
		stage.Steps = append(stage.Steps, createScriptStep(image,
			"npm_e2e",
			"npm run e2e",
		))
	}

	// add well-known e2e docker if infra is cloud
	if _, ok := json.Scripts["e2e:docker"]; ok && image == "" {
		stage.Steps = append(stage.Steps, createScriptStep(image,
			"npm_e2e_docker",
			"npm run e2e docker",
		))
	}

	// add well-known dist command
	if _, ok := json.Scripts["dist"]; ok {
		stage.Steps = append(stage.Steps, createScriptStep(image,
			"npm_dist",
			"npm run dist",
		))
	}

	return nil
}

// represents the package.json file format.
type packageJson struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Scripts map[string]interface{} `json:"scripts"`
}
