package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/MrEhbr/go-fsm"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/imports"
	"gopkg.in/yaml.v2"
)

var ActionsGenCommand = &cli.Command{
	Name:  "gen",
	Usage: "generates fsm actions",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "template",
			Aliases:  []string{"tpl"},
			Usage:    "template for action",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output_dir",
			Aliases:  []string{"o"},
			Usage:    "output dir",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "transitions",
			Aliases:  []string{"t"},
			Usage:    "path to file with transitions",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		if c.IsSet("transitions") {
			path, err := filepath.Abs(c.String("transitions"))
			if err != nil {
				return err
			}

			if err := c.Set("transitions", path); err != nil {
				return cli.Exit(err, 1)
			}
		}

		if c.IsSet("template") {
			path, err := filepath.Abs(c.String("template"))
			if err != nil {
				return err
			}

			if err := c.Set("template", path); err != nil {
				return cli.Exit(err, 1)
			}
		}

		tpl, err := ioutil.ReadFile(c.String("template"))
		if err != nil {
			return cli.Exit(err, 1)
		}

		if err := actionsGenerate(c.String("transitions"), c.String("output_dir"), string(tpl)); err != nil {
			return cli.Exit(err, 1)
		}

		return nil
	},
}

func actionsGenerate(trFile, outputDir, tpl string) error {
	transitions, err := transitionsFromFile(trFile)
	if err != nil {
		return fmt.Errorf("failed to read transitions: %w", err)
	}

	if !isDirectory(outputDir) {
		return fmt.Errorf("%s is not directory", outputDir)
	}

	funcs := map[string]interface{}{
		"to_camel": fsm.ToCamelCase,
		"to_snake": fsm.ToSnackCase,
	}

	compiledTpl, err := template.New("action").Funcs(funcs).Parse(tpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	actions := transitions.Actions()
	stubs := make(map[string][]byte, len(actions))
	for _, v := range actions {
		var buf bytes.Buffer

		if err := compiledTpl.Execute(&buf, map[string]string{"Action": v}); err != nil {
			return fmt.Errorf("failed to generate stub for %q: %w", v, err)
		}

		stubs[fsm.ToSnackCase(v, false)] = buf.Bytes()
	}

	for k, v := range stubs {
		fileName := filepath.Join(outputDir, fmt.Sprintf("%s.go", k))
		if fileExists(fileName) {
			continue
		}

		processedSource, err := imports.Process(fileName, v, nil)
		if err != nil {
			return fmt.Errorf("failed to format generated code: %w", err)
		}

		if err := ioutil.WriteFile(fileName, processedSource, 0o644); err != nil {
			return fmt.Errorf("failed to write file %q: %w", fileName, err)
		}
	}

	return nil
}

func fileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func transitionsFromFile(name string) (fsm.Transitions, error) {
	var transitions fsm.Transitions
	if name != "" {
		decoder := json.Unmarshal
		ext := filepath.Ext(name)
		if ext == ".yaml" || ext == ".yml" {
			decoder = yaml.Unmarshal
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			return nil, err
		}

		if err := decoder(data, &transitions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transitions file: %w", err)
		}
	}

	return transitions, nil
}
