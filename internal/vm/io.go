package vm

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// var declares or overwrites a variable.
// Usage: var <name> <value>
func (vm *VM) opVar(args []string, line int) error {
	if len(args) < 2 {
		return types.LineError(line, "usage: var <name> <value>")
	}
	name := args[0]
	vm.vars[name] = vm.resolve(args[1])
	return nil
}

// mov copies a value into a variable.
// Usage: mov <dest> <src>
func (vm *VM) opMov(args []string, line int) error {
	if len(args) < 2 {
		return types.LineError(line, "usage: mov <dest> <src>")
	}
	dest := args[0]
	vm.vars[dest] = vm.resolve(args[1])
	return nil
}

// print outputs a value.
// Usage: print <value>
func (vm *VM) opPrint(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: print <value>")
	}
	v := vm.resolve(args[0])
	vm.io.Print(v.Format())
	return nil
}

// input pauses execution and waits for user input.
// Usage: input <var> [prompt]
func (vm *VM) opInput(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: input <var> [prompt]")
	}
	varName := args[0]
	prompt := ""
	if len(args) > 1 {
		prompt = args[1]
	}
	vm.io.Print("__INPUT_REQUIRED__" + prompt)
	result := vm.io.Input(prompt)
	vm.vars[varName] = types.NewString(result)
	return nil
}
