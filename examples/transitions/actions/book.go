package actions

import "github.com/MrEhbr/go-fsm/v2/examples/transitions"

type BookAction struct{}

func (action *BookAction) Process(order transitions.Order) error {
	panic("not implemented")
}
