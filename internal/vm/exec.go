package vm

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// execute dispatches a single instruction to the appropriate handler.
func (vm *VM) execute(inst types.Instruction) error {
	switch inst.Op {
	// I/O
	case "var":
		return vm.opVar(inst.Args, inst.Line)
	case "mov":
		return vm.opMov(inst.Args, inst.Line)
	case "print":
		return vm.opPrint(inst.Args, inst.Line)
	case "input":
		return vm.opInput(inst.Args, inst.Line)

	// arithmetic
	case "add":
		return vm.opAdd(inst.Args, inst.Line)
	case "sub":
		return vm.opSub(inst.Args, inst.Line)
	case "mul":
		return vm.opMul(inst.Args, inst.Line)
	case "div":
		return vm.opDiv(inst.Args, inst.Line)
	case "mod":
		return vm.opMod(inst.Args, inst.Line)
	case "inc":
		return vm.opInc(inst.Args, inst.Line)
	case "dec":
		return vm.opDec(inst.Args, inst.Line)
	case "neg":
		return vm.opNeg(inst.Args, inst.Line)

	// logic
	case "and":
		return vm.opAnd(inst.Args, inst.Line)
	case "or":
		return vm.opOr(inst.Args, inst.Line)
	case "xor":
		return vm.opXor(inst.Args, inst.Line)
	case "not":
		return vm.opNot(inst.Args, inst.Line)

	// control flow
	case "label":
		return nil // already resolved during parsing
	case "jump":
		return vm.opJump(inst.Args, inst.Line)
	case "beq":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a == b })
	case "bne":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a != b })
	case "blt":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a < b })
	case "bgt":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a > b })
	case "ble":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a <= b })
	case "bge":
		return vm.opBranch(inst.Args, inst.Line, func(a, b float64) bool { return a >= b })
	case "halt":
		vm.halted = true
		return nil
	case "nop":
		return nil

	// stack
	case "push":
		return vm.opPush(inst.Args, inst.Line)
	case "pop":
		return vm.opPop(inst.Args, inst.Line)
	case "call":
		return vm.opCall(inst.Args, inst.Line)
	case "ret":
		return vm.opRet(inst.Line)

	default:
		return types.LineError(inst.Line, "unknown instruction: %s", inst.Op)
	}
}
