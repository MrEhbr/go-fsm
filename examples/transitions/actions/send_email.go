package actions

import "github.com/MrEhbr/go-fsm/v2/examples/transitions"

type SendEmailAction struct{}

func (action *SendEmailAction) Process(order transitions.Order) error {
	return nil
}
