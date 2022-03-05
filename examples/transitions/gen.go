//go:generate go-fsm actions gen --tpl ./action.go.tpl -t ./transitions.json -o ./actions
//go:generate go-fsm actions doc -t ./transitions.json -o ./actions/README.md
package transitions
