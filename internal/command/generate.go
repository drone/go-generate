package command

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/drone/go-generate/builder"
	"github.com/drone/go-generate/utils/chroot"
	"github.com/drone/go-generate/utils/cloner"
	"github.com/google/subcommands"
)

type Generate struct {
	username   string
	password   string
	privatekey string
}

func (*Generate) Name() string     { return "generate" }
func (*Generate) Synopsis() string { return "generate generates a pipeline" }
func (*Generate) Usage() string {
	return `generate [-username] [-password] <repository>
`
}

func (c *Generate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.username, "username", "", "repository username")
	f.StringVar(&c.password, "password", "", "repository password")
	f.StringVar(&c.privatekey, "privatekey", "", "repositroy private key")
}

func (c *Generate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// extract the repository path
	path := f.Arg(0)

	// if the path is a repository url,
	// clone the rempository into a temporary
	// directory, and then delete
	if isRemote(path) {
		temp, err := os.MkdirTemp("", "")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return subcommands.ExitFailure
		}
		defer os.RemoveAll(temp)

		params := cloner.Params{
			Dir:        temp,
			Repo:       path,
			Username:   "", // not yet implemented
			Password:   "", // not yet implemented
			Privatekey: "", // not yet implemented
		}
		cloner := cloner.New(1, io.Discard) // 1 depth, discard git clone logs
		cloner.Clone(context.Background(), params)

		// change the path to the temp directory
		path = temp
	}

	// create a chroot virtual filesystem that we
	// pass to the builder for isolation purposes.
	chroot, err := chroot.New(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	// builds the pipeline configuration based on
	// the contents of the virtual filesystem.
	builder := builder.New()
	out, err := builder.Build(chroot)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	// output to console
	os.Stdout.Write(out)

	return subcommands.ExitSuccess
}

// returns true if the string is a remote git repository.
func isRemote(s string) bool {
	return strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@")
}
