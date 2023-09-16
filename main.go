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

package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/drone/go-generate/builder"
	"github.com/drone/go-generate/chroot"
	"github.com/drone/go-generate/cloner"
)

func main() {
	var path string

	// extract the repository path
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// if the path is a repository url,
	// clone the rempository into a temporary
	// directory, and then delete
	if isRemote(path) {
		temp, err := ioutil.TempDir("", "")
		if err != nil {
			log.Fatalln(err)
		}
		defer os.RemoveAll(temp)

		params := cloner.Params{
			Dir:        temp,
			Repo:       path,
			Username:   "", // not yet implemented
			Password:   "", // not yet implemented
			Privatekey: "", // not yet implemented
		}
		cloner := cloner.New(1, ioutil.Discard) // 1 depth, discard git clone logs
		cloner.Clone(context.Background(), params)

		// change the path to the temp directory
		path = temp
	}

	// create a chroot virtual filesystem that we
	// pass to the builder for isolation purposes.
	chroot, err := chroot.New(path)
	if err != nil {
		log.Fatalln(err)
	}

	// builds the pipeline configuration based on
	// the contents of the virtual filesystem.
	builder := builder.New()
	out, err := builder.Build(chroot)
	if err != nil {
		log.Fatalln(err)
	}

	// output to console
	println(string(out))
}

// returns true if the string is a remote git repository.
func isRemote(s string) bool {
	return strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@")
}
