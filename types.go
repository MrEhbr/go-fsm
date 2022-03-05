package fsm

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"moul.io/graphman"
)

// Struct represent struct data
type Struct struct {
	Name        string
	StateField  string
	StateType   string
	StateValues []StateValue
	Transitions Transitions
}

// StateValue represents state field values
type StateValue struct {
	Name string
	Val  string
}

// Validate validate Struct for using in generation
func (s Struct) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("struct name is empty")
	}

	if s.StateField == "" {
		return fmt.Errorf("state field is empty")
	}

	if s.StateType == "" {
		return fmt.Errorf("state type is empty")
	}

	if len(s.Transitions) > 0 && len(s.StateValues) == 0 {
		return fmt.Errorf("no state values for transitions")
	}

	for _, t := range s.Transitions {
		for _, f := range t.From {
			if v := s.FindValue(f); v == "" {
				return fmt.Errorf("unknown state value[%s]", f)
			}
		}

		if v := s.FindValue(t.To); v == "" {
			return fmt.Errorf("unknown state value[%s]", t.To)
		}
	}

	return nil
}

// FindValue find state value by string representation of value or variable name
// if not found return empty string
func (s Struct) FindValue(str string) string {
	if strings.TrimSpace(str) == "" {
		return ""
	}

	for _, v := range s.StateValues {
		if strings.EqualFold(v.Name, str) || v.Val == str {
			return v.Name
		}

		if strings.EqualFold(ToCamelCase(v.Name, true), ToCamelCase(str, true)) {
			return v.Name
		}

		trimmed := strings.TrimPrefix(v.Name, s.StateType)
		if strings.EqualFold(trimmed, str) {
			return v.Name
		}

		if strings.EqualFold(ToCamelCase(trimmed, true), ToCamelCase(str, true)) {
			return v.Name
		}
	}

	return ""
}

// Transition represent fsm state transition
type Transition struct {
	From          stringArray `json:"from" yaml:"from"`
	To            string      `json:"to" yaml:"to"`
	Event         string      `json:"event" yaml:"event"`
	BeforeActions []string    `json:"before_actions" yaml:"before_actions"`
	Actions       []string    `json:"actions" yaml:"actions"`
}

func (t Transition) String() string {
	var builder strings.Builder

	for _, v := range t.From {
		if builder.Len() > 0 {
			builder.WriteRune('\n')
		}
		builder.WriteRune('-')
		builder.WriteRune(' ')
		builder.WriteString(v)
		builder.WriteString("->")
		builder.WriteString(t.To)
	}

	return builder.String()
}

type stringArray []string

func (a *stringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

// Transitions slice of Transition
type Transitions []Transition

// Actions return unique slice of actions
func (trs Transitions) Actions() []string {
	if len(trs) == 0 {
		return nil
	}

	actions := make([]string, 0, 10)
	for _, v := range trs {
		actions = append(actions, v.Actions...)
		actions = append(actions, v.BeforeActions...)
	}

	for i := range actions {
		actions[i] = filterActionName(actions[i])
	}

	sort.Strings(actions)
	j := 0
	for i := 1; i < len(actions); i++ {
		if actions[j] == actions[i] {
			continue
		}
		j++
		// preserve the original data
		// actions[i], actions[j] = actions[j], actions[i]
		// only set what is required
		actions[j] = actions[i]
	}

	return actions[:j+1]
}

func (trs Transitions) Events() []string {
	if len(trs) == 0 {
		return nil
	}

	events := make([]string, 0, 10)
	for _, v := range trs {
		if v.Event == "" {
			continue
		}

		events = append(events, v.Event)
	}

	sort.Strings(events)
	j := 0
	for i := 1; i < len(events); i++ {
		if events[j] == events[i] {
			continue
		}
		j++
		// preserve the original data
		// events[i], events[j] = events[j], events[i]
		// only set what is required
		events[j] = events[i]
	}

	return events[:j+1]
}

func (trs Transitions) Graph() *graphman.Graph {
	graph := graphman.New(graphman.Attrs{
		"overlap": "scalexy",
		"splines": "true",
	})

	for _, tr := range trs {
		attrs := graphman.Attrs{}
		attrs.SetTitle(tr.Event)

		for _, from := range tr.From {
			graph.AddEdge(from, tr.To, attrs)
		}
	}

	colors := NewColors()
	for _, v := range graph.Vertices() {
		v.Attrs.SetPertState()

		if v.IsSink() || v.IsSource() {
			v.SetColor(colors.Pick())
			v.Attrs["style"] = "bold"
		}
	}

	return graph
}

func (trs Transitions) ActionTransitions(action string) []Transition {
	res := make([]Transition, 0)
	for _, v := range trs {
		for _, b := range v.BeforeActions {
			if b == action {
				res = append(res, v)
			}
			continue
		}
		for _, a := range v.Actions {
			if a == action {
				res = append(res, v)
			}
			continue
		}
	}

	return res
}

func filterActionName(action string) string {
	action = strings.TrimLeftFunc(action, func(r rune) bool {
		if unicode.IsLetter(r) || r == '_' {
			return false
		}

		return true
	})

	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}

		return -1
	}, action)
}
