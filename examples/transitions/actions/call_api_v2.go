package actions

import "github.com/MrEhbr/go-fsm/v2/examples/transitions"

type CallApiV2Action struct{}

func (action *CallApiV2Action) Process(order transitions.Order) error {
	panic("not implemented")
}
