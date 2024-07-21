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
	"bytes"
	"io/fs"

	spec "github.com/bradrydzewski/spec/yaml"
)

// ConfigureRails configures a Ruby on Rails step.
func ConfigureRails(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0]

	// check if we should use a container-based
	// execution environment.
	var image string
	if isContainerRuntime(pipeline) {
		image = "ruby"
	}

	// check for a ruby gemfile
	if exists(fsys, "Gemfile") {

		// ignore gemfiles that do not contain the
		// rails dependency
		gemfile, _ := read(fsys, "Gemfile")
		if !bytes.Contains(gemfile, []byte("'rails'")) {
			return nil
		}

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_install",
			"bundle install — jobs=3 — retry=3",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_db_create",
			"bundle exec rake db:create",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_db_migrate",
			"bundle exec rake db:migrate",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_rspec",
			"bundle exec rspec",
		))
	}

	return nil
}
