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

// Package builder builds a pipeline configuration.
package builder

import (
	"errors"
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/ghodss/yaml"
)

// Rule defines a pipeline build rule.
type Rule func(workspace fs.FS, pipeline *spec.Pipeline) error

// SkipAll is used as a return value from Rule to indicate
// that all remaining rules are to be skipped. It is never
// returned as an error by the Builder.
var SkipAll = errors.New("skip everything and stop the pipeline generation")

// Builder builds a pipeline configuration.
type Builder struct {
	rules []Rule
}

// New creates a new pipeline builder.
func New() *Builder {
	return &Builder{
		rules: []Rule{
			ConfigurePlatform,
			ConfigureGo,
			ConfigureNode,
			ConfigurePython,
			ConfigureRails,
			ConfigureRuby,
			ConfigureRust,
			ConfigureSwift,
			ConfigureDocker,
			ConfigureDefault,
		},
	}
}

// New creates a new pipeline builder with custom rules.
func NewRules(rules []Rule) *Builder {
	return &Builder{
		rules: rules,
	}
}

// Build the pipeline configuration.
func (b *Builder) Build(fsys fs.FS) ([]byte, error) {
	stageci := new(spec.StageCI)
	stageci.Platform = new(spec.Platform)
	stageci.Platform.Os = spec.OSLinux
	stageci.Platform.Arch = spec.ArchAmd64

	stage := new(spec.Stage)
	stage.Name = "build"
	stage.Type = "ci"
	stage.Spec = stageci

	pipeline := new(spec.Pipeline)
	pipeline.Stages = append(pipeline.Stages, stage)
	for _, rule := range b.rules {
		if err := rule(fsys, pipeline); err == SkipAll {
			break
		}

		// we purposefully ignore errors here.
		// an error in an individual rule should
		// never prevent yaml generation.
	}

	config := new(spec.Config)
	config.Kind = "pipeline"
	config.Spec = pipeline
	return yaml.Marshal(config)
}

//
// helper functions.
//

// helper function to create a script step.
func createScriptStep(image, name, command string) *spec.Step {
	script := new(spec.StepExec)
	script.Run = command

	if image != "" {
		script.Image = image
	}

	step := new(spec.Step)
	step.Name = name
	step.Type = "script"
	step.Spec = script

	return step
}
