package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drone/go-generate/updater"
	"github.com/drone/go-generate/utils/openai"
	"github.com/google/subcommands"
)

type Update struct {
	config string
	repo   string
	token  string
}

func (*Update) Name() string     { return "update" }
func (*Update) Synopsis() string { return "update updates a pipeline" }
func (*Update) Usage() string {
	return `optimize [-config] [-token] [-repo] <prompt>
`
}

func (c *Update) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.config, "config", "", "config is the path to the pipeline yaml")
	f.StringVar(&c.repo, "repo", "", "repo is the path or url to the repository")
	f.StringVar(&c.token, "token", "", "error is the openai token")
}

func (c *Update) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	in := &updater.Input{
		Prompt: f.Arg(0),
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

	// create the updater
	u := updater.New(client)

	// update the pipeline configuration
	out, err := u.Update(in)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	// TODO print output in a user-friendly way
	// once we know the output attributes
	json.NewEncoder(os.Stdout).Encode(out)

	return subcommands.ExitSuccess
}
