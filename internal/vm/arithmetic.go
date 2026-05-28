package vm

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// arithmeticOp is a helper for three-operand arithmetic instructions.
// If both operands are integers, the result is integer; otherwise float.
// Exception: div always produces float.
func (vm *VM) arithmeticOp(args []string, line int, op string) error {
	if len(args) < 3 {
		return types.LineError(line, "usage: %s <dest> <a> <b>", op)
	}
	dest := args[0]
	a := vm.resolve(args[1])
	b := vm.resolve(args[2])

	// string concatenation with add
	if op == "add" && (a.Type == types.TypeString || b.Type == types.TypeString) {
		vm.vars[dest] = types.NewString(a.Format() + b.Format())
		return nil
	}

	af := a.AsNumber()
	bf := b.AsNumber()

	var result float64
	switch op {
	case "add":
		result = af + bf
	case "sub":
		result = af - bf
	case "mul":
		result = af * bf
	case "div":
		if bf == 0 {
			return types.LineError(line, "division by zero")
		}
		result = af / bf
	case "mod":
		if bf == 0 {
			return types.LineError(line, "division by zero")
		}
		result = float64(int64(af) % int64(bf))
	}

	// use integer type when both operands are int and op is not div
	if a.Type == types.TypeInt && b.Type == types.TypeInt && op != "div" {
		vm.vars[dest] = types.NewInt(int64(result))
	} else {
		vm.vars[dest] = types.NewFloat(result)
	}
	return nil
}

func (vm *VM) opAdd(args []string, line int) error { return vm.arithmeticOp(args, line, "add") }
func (vm *VM) opSub(args []string, line int) error { return vm.arithmeticOp(args, line, "sub") }
func (vm *VM) opMul(args []string, line int) error { return vm.arithmeticOp(args, line, "mul") }
func (vm *VM) opDiv(args []string, line int) error { return vm.arithmeticOp(args, line, "div") }
func (vm *VM) opMod(args []string, line int) error { return vm.arithmeticOp(args, line, "mod") }

// inc increments a variable by 1.
func (vm *VM) opInc(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: inc <var>")
	}
	v, err := vm.getVar(args[0], line)
	if err != nil {
		return err
	}
	if v.Type == types.TypeInt {
		vm.vars[args[0]] = types.NewInt(v.Int + 1)
	} else {
		vm.vars[args[0]] = types.NewFloat(v.AsNumber() + 1)
	}
	return nil
}

// dec decrements a variable by 1.
func (vm *VM) opDec(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: dec <var>")
	}
	v, err := vm.getVar(args[0], line)
	if err != nil {
		return err
	}
	if v.Type == types.TypeInt {
		vm.vars[args[0]] = types.NewInt(v.Int - 1)
	} else {
		vm.vars[args[0]] = types.NewFloat(v.AsNumber() - 1)
	}
	return nil
}

// neg negates a numeric variable.
func (vm *VM) opNeg(args []string, line int) error {
	if len(args) < 1 {
		return types.LineError(line, "usage: neg <var>")
	}
	v, err := vm.getVar(args[0], line)
	if err != nil {
		return err
	}
	if v.Type == types.TypeInt {
		vm.vars[args[0]] = types.NewInt(-v.Int)
	} else {
		vm.vars[args[0]] = types.NewFloat(-v.AsNumber())
	}
	return nil
}
