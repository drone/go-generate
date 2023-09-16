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

import (
	"regexp"
	"strings"
)

// regular expressions to test whether or not a string is
// a sha1 or sha256 commit hash.
var (
	sha1   = regexp.MustCompile("^([a-f0-9]{40})$")
	sha256 = regexp.MustCompile("^([a-f0-9]{64})$")
)

// helper function returns true if the string is a commit hash.
func isHash(s string) bool {
	return sha1.MatchString(s) || sha256.MatchString(s)
}

// helper function returns the branch name expanded to the
// fully qualified reference path (e.g refs/heads/master).
func expandRef(name string) string {
	if strings.HasPrefix(name, "refs/") {
		return name
	}
	return "refs/heads/" + name
}
