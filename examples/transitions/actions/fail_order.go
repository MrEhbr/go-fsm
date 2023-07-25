package actions

import "github.com/MrEhbr/go-fsm/v2/examples/transitions"

type FailOrderAction struct{}

func (action *FailOrderAction) Process(order transitions.Order) error {
	panic("not implemented")
}
