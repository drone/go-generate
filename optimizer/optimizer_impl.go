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

package optimizer

import (
	"github.com/drone/go-generate/utils/openai"
)

type OptimizerImpl struct {
	client openai.Client
}

func New(client openai.Client) *OptimizerImpl {
	return &OptimizerImpl{client: client}
}

func (d *OptimizerImpl) Optimize(in *Input) (*Output, error) {

	// FIXME: for now we just echo the yaml input as the output
	out := new(Output)
	out.Before = in.Config
	out.After = in.Config

	return out, nil
}