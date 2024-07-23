package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/drone/go-generate/debugger"
	"github.com/drone/go-generate/utils/openai"

	"github.com/google/subcommands"
)

type Debug struct {
	config string
	repo   string
	logs   string
	error  string
	token  string
}

func (*Debug) Name() string     { return "debug" }
func (*Debug) Synopsis() string { return "debug debugs a failed pipeline" }
func (*Debug) Usage() string {
	return `debug [-config] [-logs] [-error] [-token] [-repo]
`
}

func (c *Debug) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.config, "config", "", "config is the path to the pipeline yaml")
	f.StringVar(&c.logs, "logs", "", "logs is the path to the pipeline log file")
	f.StringVar(&c.error, "error", "", "error is the path the pipeline error message")
	f.StringVar(&c.repo, "repo", "", "repo is the path or url to the repository")
	f.StringVar(&c.token, "token", "", "error is the openai token")
}

func (c *Debug) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// path := f.Arg(0)

	in := &debugger.Input{
		Config: nil,
		Logs:   nil,
		Error:  nil,
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

	// read the logs
	if path := c.logs; path != "" {
		out, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return subcommands.ExitFailure
		}
		in.Logs = out
	}

	// read the error
	if path := c.error; path != "" {
		out, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return subcommands.ExitFailure
		}
		in.Error = out
	}

	// TODO clone the repository and set in.Code
	// as the virtual file system. See generate.go:42

	// create the openai client
	client := openai.New(c.token)

	// create the debugger
	d := debugger.New(client)

	// debug the failure
	out, err := d.Debug(in)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	// TODO print output in a user-friendly way
	// once we know the output attributes
	json.NewEncoder(os.Stdout).Encode(out)

	return subcommands.ExitSuccess
}
