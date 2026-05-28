package vm

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// push pushes a value onto the data stack.
func (vm *VM) opPush(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: push <value>")
	}
	v := vm.resolve(args[0])
	vm.stack = append(vm.stack, v)
	return nil
}

// pop pops the top value from the data stack into a variable.
func (vm *VM) opPop(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: pop <var>")
	}
	if len(vm.stack) == 0 {
		return types.LineError(line, "stack underflow")
	}
	top := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	vm.vars[args[0]] = top
	return nil
}

// call pushes the return address onto the call stack and jumps to a label.
func (vm *VM) opCall(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: call <target>")
	}
	idx, err := vm.resolveLabel(args[0], line)
	if err != nil {
		return err
	}
	vm.calls = append(vm.calls, vm.pc)
	vm.pc = idx
	return nil
}

// ret pops the return address from the call stack and jumps back.
func (vm *VM) opRet(line int) error {
	if len(vm.calls) == 0 {
		return types.LineError(line, "call stack underflow")
	}
	retAddr := vm.calls[len(vm.calls)-1]
	vm.calls = vm.calls[:len(vm.calls)-1]
	vm.pc = retAddr
	return nil
}
