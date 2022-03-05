# go-fsm

![CI](https://github.com/MrEhbr/go-fsm/workflows/CI/badge.svg)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/MrEhbr/go-fsm/blob/main/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/MrEhbr/go-fsm.svg)](https://github.com/MrEhbr/go-fsm/releases)
[![codecov](https://codecov.io/gh/MrEhbr/go-fsm/branch/main/graph/badge.svg)](https://codecov.io/gh/MrEhbr/go-fsm)
![Made by Alexey Burmistrov](https://img.shields.io/badge/made%20by-Alexey%20Burmistrov-blue.svg?style=flat)

go-fsm is a command line tool that generates finite state machine for Go struct.

## Usage

### FSM generation

```console
Usage: go-fsm gen -p ./examples/transitions -s Order -f State -o order_fsm.go -t ./examples/transitions/transitions.json
   --package, -p      package where struct is located (default: default is current dir(.))
   --struct, -s       struct name
   --field, -f        state field of struct
   --output, -o       output file name (default: default srcdir/<struct>_fsm.go)
   --transitions, -t  path to file with transitions
   --noGenerate, -g   don't put //go:generate instruction to the generated code (default: false)
   --graph-output, -a value  path to transition graph file in dot format
```

This will generate [finite state machine](./examples/transitions/order_fsm.go) for struct Order with transitions defined in [./examples/transitions/transitions.json](./examples/transitions/transitions.json) file.
Transition graph will be generated from transitions file and output will be in dot format

### Action generation

```console
Usage: go-fsm actions gen --tpl examples/transitions/action.go.tpl -o examples/transitions/actions -t examples/transitions/transitions.json
   --template value, --tpl value  template for action
   --output_dir value, -o value   output dir
   --transitions value, -t value  path to file with transitions
```

This will generate [actions](./examples/transitions/actions) stubs from [./examples/transitions/action.go.tpl](./examples/transitions/action.go.tpl) file.

### Action doc generation

```console
Usage: go-fsm actions doc -o examples/transitions/actions/README.md -t examples/transitions/transitions.json
   --output value, -o value       output file name
   --transitions value, -t value  path to file with transitions
```

This will generate actions doc [README.md](./examples/transitions/actions/README.md) from [transitions](./examples/transitions/transitions.json) file.
Further runs of the command will add only new actions and update the `Transitions where action appears` section.

## Install

### Using go

```console
go get -u github.com/MrEhbr/go-fsm/cmd/go-fsm
```

### Download releases

<https://github.com/MrEhbr/go-fsm/releases>

## License

© 2020 [Alexey Burmistrov]

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE`](LICENSE)). See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: Apache-2.0`
