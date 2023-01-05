package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/MrEhbr/go-fsm"
	"github.com/MrEhbr/go-fsm/cmd/go-fsm/commands"
	"github.com/matryer/is"
)

func TestActionsDocCommand(t *testing.T) {
	const trsFile = "./test_data/transitions.json"
	var transitions fsm.Transitions
	data, err := os.ReadFile(trsFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &transitions); err != nil {
		t.Fatal("failed to unmarshal transitions: %w", err)
	}

	t.Run("first time generate", func(t *testing.T) {
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile, err := os.CreateTemp(tmpDir, "README.md")
		is.NoErr(err)
		is.NoErr(commands.ActionDocGenerate(trsFile, tmpFile.Name()))
		actions := transitions.Actions()
		doc, err := io.ReadAll(tmpFile)
		is.NoErr(err)

		for _, v := range actions {
			is.True(strings.Contains(string(doc), fmt.Sprintf("### %s", v)))
			is.True(strings.Contains(string(doc), fmt.Sprintf("#### Transitions where action %q appears", v)))
			is.True(strings.Contains(string(doc), fmt.Sprintf("#### %s description", v)))
		}
	})

	t.Run("append actions, transitions updated", func(t *testing.T) {
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile, err := os.CreateTemp(tmpDir, "README.md")
		is.NoErr(err)

		want := `### book

#### Transitions where action "book" appears

- CREATED->STARTED

#### book description

some description

---`
		tpl := `### book

#### Transitions where action "book" appears

- UNKNOWN->UNKNOWN

#### book description

some description

---`
		is.NoErr(os.WriteFile(tmpFile.Name(), []byte(tpl), 0o644))
		is.NoErr(commands.ActionDocGenerate(trsFile, tmpFile.Name()))
		actions := transitions.Actions()
		doc, err := io.ReadAll(tmpFile)
		is.NoErr(err)

		is.True(strings.Contains(string(doc), want))
		for _, v := range actions {
			is.True(strings.Contains(string(doc), fmt.Sprintf("### %s", v)))
			is.True(strings.Contains(string(doc), fmt.Sprintf("#### Transitions where action %q appears", v)))
			is.True(strings.Contains(string(doc), fmt.Sprintf("#### %s description", v)))
		}
	})
}
