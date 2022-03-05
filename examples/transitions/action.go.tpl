package actions

type {{to_camel .Action true}}Action struct {

}

func(action *{{to_camel .Action true}}Action) Process(order transitions.Order) error {
    panic("not implemented")
}
