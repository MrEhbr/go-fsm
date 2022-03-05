package actions

import "github.com/MrEhbr/go-fsm/examples/transitions"

type BookAction struct{}

func (action *BookAction) Process(order transitions.Order) error {
	panic("not implemented")
}
