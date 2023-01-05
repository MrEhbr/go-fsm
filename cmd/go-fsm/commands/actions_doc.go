package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kunde21/markdownfmt/v2/markdown"
	"github.com/Kunde21/markdownfmt/v2/markdownfmt"
	"github.com/MrEhbr/go-fsm"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var ActionsDocCommand = &cli.Command{
	Name:  "doc",
	Usage: "generates action doc",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "output file name",
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

		if err := ActionDocGenerate(c.String("transitions"), c.String("output")); err != nil {
			return cli.Exit(err, 1)
		}

		return nil
	},
}

const (
	actionDocTpl = `
### %s

#### Transitions where action %q appears

%s

#### %s description

---

`
)

func ActionDocGenerate(trFile, output string) error {
	transitions, err := transitionsFromFile(trFile)
	if err != nil {
		return fmt.Errorf("failed to read transitions: %w", err)
	}

	if isDirectory(output) {
		return fmt.Errorf("%s is directory", output)
	}

	var source []byte
	source, _ = os.ReadFile(output)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRenderer(markdown.NewRenderer()),
		goldmark.WithParserOptions(
			parser.WithHeadingAttribute(),
		),
	)

	doc := md.Parser().Parse(text.NewReader(source))

	actions := transitions.Actions()
	haveDocs := make(map[string]bool, len(actions))

	var actionName string
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkSkipChildren, nil
		}

		// we interested only in headings
		if n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}
		// in level 3 heading written action name
		if n.(*ast.Heading).Level == 3 {
			actionName = string(n.Text(source))
			haveDocs[actionName] = true
			return ast.WalkContinue, nil
		}

		if n.(*ast.Heading).Level == 4 && strings.HasPrefix(strings.ToLower(string(n.Text(source))), "transitions") {
			next := n.NextSibling()
			paragraph := createParagraph(transitions, actionName)
			if paragraph == nil {
				return ast.WalkContinue, nil
			}
			if next != nil && next.Kind() == ast.KindList {
				if next.ChildCount() > 0 {
					next.RemoveChildren(next)
				}
				next.AppendChild(next, paragraph)
			} else {
				n.AppendChild(n, paragraph)
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return fmt.Errorf("markdown traversal error: %w", err)
	}

	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, source, doc)
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	for _, v := range actions {
		if haveDocs[v] {
			continue
		}

		actionDocStr := fmt.Sprintf(actionDocTpl, v, v, actionTransitions(transitions, v), v)
		buf.WriteString(actionDocStr)
	}

	data, err := markdownfmt.Process(output, buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format markdown: %w", err)
	}

	return os.WriteFile(output, data, 0o644)
}

func actionTransitions(trs fsm.Transitions, action string) string {
	actionTrs := trs.ActionTransitions(action)
	var builder strings.Builder
	for i, v := range actionTrs {
		builder.WriteString(v.String())
		if i != len(actionTrs)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}

func createParagraph(trs fsm.Transitions, action string) *ast.Paragraph {
	doc := actionTransitions(trs, action)
	if doc == "" {
		return nil
	}
	txt := ast.NewText()
	txt.AppendChild(txt, ast.NewString([]byte(doc)))
	txt.SetRaw(true)
	textBlock := ast.NewTextBlock()
	textBlock.AppendChild(textBlock, txt)
	paragraph := ast.NewParagraph()
	paragraph.AppendChild(paragraph, textBlock)
	return paragraph
}
