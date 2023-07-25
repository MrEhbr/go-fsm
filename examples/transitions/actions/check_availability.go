package actions

import "github.com/MrEhbr/go-fsm/v2/examples/transitions"

type CheckAvailabilityAction struct{}

func (action *CheckAvailabilityAction) Process(order transitions.Order) error {
	panic("not implemented")
}
