package commands

import (
	"github.com/urfave/cli/v2"
)

var ActionsCommand = &cli.Command{
	Name:  "actions",
	Usage: "fsm actions",
	Subcommands: []*cli.Command{
		ActionsGenCommand,
		ActionsDocCommand,
	},
}
