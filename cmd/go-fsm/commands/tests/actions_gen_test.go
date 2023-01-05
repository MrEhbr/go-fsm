package tests

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MrEhbr/go-fsm"
	"github.com/MrEhbr/go-fsm/cmd/go-fsm/commands"
	"github.com/matryer/is"
	"github.com/urfave/cli/v2"
)

func TestActionsGenCommand(t *testing.T) {
	t.Run("required options", func(t *testing.T) {
		is := is.New(t)

		app := &cli.App{
			Writer:         io.Discard,
			ExitErrHandler: func(_ *cli.Context, _ error) {},
		}
		set := flag.NewFlagSet("test", 0)
		is.NoErr(set.Parse([]string{
			commands.ActionsGenCommand.Name,
		}))
		c := cli.NewContext(app, set, nil)
		err := commands.ActionsGenCommand.Run(c)
		is.True(err != nil)                                                 // must be error
		is.True(strings.Contains(strings.ToLower(err.Error()), "required")) // error must contains required msg
		is.True(err != nil)                                                 // wanted error with required flags
	})

	t.Run("template not found", func(t *testing.T) {
		is := is.New(t)

		tempDir := t.TempDir()
		app := &cli.App{
			Writer:         io.Discard,
			ExitErrHandler: func(_ *cli.Context, _ error) {},
		}
		set := flag.NewFlagSet("test", 0)
		is.NoErr(set.Parse([]string{
			commands.ActionsGenCommand.Name,
			"-template", "not_exists",
			"-output_dir", tempDir,
			"-transitions", "./test_data/transitions.json",
		}))
		c := cli.NewContext(app, set, nil)
		err := commands.ActionsGenCommand.Run(c)
		is.True(err != nil) // wanted error with required flags
		is.True(strings.Contains(strings.ToLower(err.Error()), "no such file"))
	})

	t.Run("transitions not found", func(t *testing.T) {
		is := is.New(t)

		tempDir := t.TempDir()
		app := &cli.App{
			Writer:         io.Discard,
			ExitErrHandler: func(_ *cli.Context, _ error) {},
		}
		set := flag.NewFlagSet("test", 0)
		is.NoErr(set.Parse([]string{
			commands.ActionsGenCommand.Name,
			"-template", "./test_data/action.go.tpl",
			"-output_dir", tempDir,
			"-transitions", "./test_data/not_exists.json",
		}))
		c := cli.NewContext(app, set, nil)
		err := commands.ActionsGenCommand.Run(c)
		is.True(err != nil) // wanted error with required flags
		is.True(strings.Contains(strings.ToLower(err.Error()), "no such file"))
	})

	t.Run("all actions stubs were generated", func(t *testing.T) {
		is := is.New(t)

		tempDir := t.TempDir()
		app := &cli.App{
			Writer:         io.Discard,
			ExitErrHandler: func(_ *cli.Context, _ error) {},
		}
		set := flag.NewFlagSet("test", 0)
		is.NoErr(set.Parse([]string{
			commands.ActionsGenCommand.Name,
			"-template", "./test_data/action.go.tpl",
			"-output_dir", tempDir,
			"-transitions", "./test_data/transitions.json",
		}))
		c := cli.NewContext(app, set, nil)
		err := commands.ActionsGenCommand.Run(c)
		is.NoErr(err)
		files, err := os.ReadDir(tempDir)
		is.NoErr(err)
		trsData, err := os.ReadFile("./test_data/transitions.json")
		is.NoErr(err)
		var trs fsm.Transitions
		is.NoErr(json.Unmarshal(trsData, &trs))
		actions := trs.Actions()
		is.True(len(actions) == len(files))
		for _, v := range files {
			is.True(filepath.Ext(v.Name()) == ".go")
			info, _ := v.Info()
			is.True(!v.IsDir() && info.Size() > 0)
			is.True(fsm.InStrings(actions, strings.TrimSuffix(v.Name(), ".go")))
		}
	})
}
