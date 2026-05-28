package vm

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// jump unconditionally jumps to a label.
func (vm *VM) opJump(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: jump <target>")
	}
	idx, err := vm.resolveLabel(args[0], line)
	if err != nil {
		return err
	}
	vm.pc = idx
	return nil
}

// opBranch handles all conditional branches (beq, bne, blt, bgt, ble, bge).
func (vm *VM) opBranch(args []string, line int, cond func(a, b float64) bool) error {
	if len(args) < 3 {
		return types.LineError(line, "usage: bxx <a> <b> <target>")
	}
	a := vm.resolve(args[0])
	b := vm.resolve(args[1])
	idx, err := vm.resolveLabel(args[2], line)
	if err != nil {
		return err
	}
	if cond(a.AsNumber(), b.AsNumber()) {
		vm.pc = idx
	}
	return nil
}

// logic operations

func (vm *VM) opAnd(args []string, line int) error {
	if len(args) < 3 {
		return types.LineError(line, "usage: and <dest> <a> <b>")
	}
	a := vm.resolve(args[1])
	b := vm.resolve(args[2])
	vm.vars[args[0]] = types.NewBool(a.AsBool() && b.AsBool())
	return nil
}

func (vm *VM) opOr(args []string, line int) error {
	if len(args) < 3 {
		return types.LineError(line, "usage: or <dest> <a> <b>")
	}
	a := vm.resolve(args[1])
	b := vm.resolve(args[2])
	vm.vars[args[0]] = types.NewBool(a.AsBool() || b.AsBool())
	return nil
}

func (vm *VM) opXor(args []string, line int) error {
	if len(args) < 3 {
		return types.LineError(line, "usage: xor <dest> <a> <b>")
	}
	a := vm.resolve(args[1])
	b := vm.resolve(args[2])
	vm.vars[args[0]] = types.NewBool(a.AsBool() != b.AsBool())
	return nil
}

func (vm *VM) opNot(args []string, line int) error {
	if len(args) < 2 {
		return types.LineError(line, "usage: not <dest> <a>")
	}
	a := vm.resolve(args[1])
	vm.vars[args[0]] = types.NewBool(!a.AsBool())
	return nil
}
