package actions

type {{to_camel .Action true}}Action struct {

}

func(action *{{to_camel .Action true}}Action) Process() error {
    panic("not implemented")
}
