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

package cloner

import "testing"

func TestExpandRef(t *testing.T) {
	tests := []struct {
		name, prefix, after string
	}{
		// branch references
		{
			after:  "refs/heads/master",
			name:   "master",
			prefix: "refs/heads",
		},
		{
			after:  "refs/heads/master",
			name:   "master",
			prefix: "refs/heads/",
		},
		// is already a ref
		{
			after:  "refs/tags/v1.0.0",
			name:   "refs/tags/v1.0.0",
			prefix: "refs/heads/",
		},
	}
	for _, test := range tests {
		if got, want := expandRef(test.name), test.after; got != want {
			t.Errorf("Got reference %s, want %s", got, want)
		}
	}
}

func TestIsHash(t *testing.T) {
	tests := []struct {
		name string
		tag  bool
	}{
		{
			name: "aacad6eca956c3a340ae5cd5856aa9c4a3755408",
			tag:  true,
		},
		{
			name: "3da541559918a808c2402bba5012f6c60b27661c",
			tag:  true,
		},
		{
			name: "f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b",
			tag:  true,
		},
		// not a sha
		{
			name: "aacad6e",
			tag:  false,
		},
		{
			name: "master",
			tag:  false,
		},
		{
			name: "refs/heads/master",
			tag:  false,
		},
		{
			name: "issue/42",
			tag:  false,
		},
		{
			name: "feature/foo",
			tag:  false,
		},
	}
	for _, test := range tests {
		if got, want := isHash(test.name), test.tag; got != want {
			t.Errorf("Detected hash %v, want %v", got, want)
		}
	}
}
