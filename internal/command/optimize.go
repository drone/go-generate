package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drone/go-generate/optimizer"
	"github.com/drone/go-generate/utils/openai"
	"github.com/google/subcommands"
)

type Optimize struct {
	config string
	repo   string
	token  string
}

func (*Optimize) Name() string     { return "optimize" }
func (*Optimize) Synopsis() string { return "optimize optimizes a pipeline" }
func (*Optimize) Usage() string {
	return `optimize [-config] [-token] [-repo]
`
}

func (c *Optimize) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.config, "config", "", "config is the path to the pipeline yaml")
	f.StringVar(&c.repo, "repo", "", "repo is the path or url to the repository")
	f.StringVar(&c.token, "token", "", "error is the openai token")
}

func (c *Optimize) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// user can optionally provide yaml path
	// as an arg instead of a flag.
	if path := f.Arg(0); path != "" {
		c.config = path
	}

	in := &optimizer.Input{
		Config: nil,
		Code:   nil,
	}

	// read the config
	if path := c.config; path != "" {
		out, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return subcommands.ExitFailure
		}
		in.Config = out
	}

	// TODO clone the repository and set in.Code
	// as the virtual file system. See generate.go:42

	// create the openai client
	client := openai.New(c.token)

	// create the optimizer
	o := optimizer.New(client)

	// optimize the pipeline configuration
	out, err := o.Optimize(in)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	// TODO print output in a user-friendly way
	// once we know the output attributes
	json.NewEncoder(os.Stdout).Encode(out)

	return subcommands.ExitSuccess
}
