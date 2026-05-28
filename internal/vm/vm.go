package vm

import (
	"strconv"

	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// DefaultMaxSteps is the default execution step limit to prevent infinite loops.
const DefaultMaxSteps = 1_000_000

// VM is the pigsh virtual machine that executes a parsed Program.
type VM struct {
	program  *types.Program
	pc       int // program counter (index into program.Instructions)
	vars     map[string]types.Value
	stack    []types.Value
	calls    []int // call stack (return addresses)
	io       types.IO
	halted   bool
	maxSteps int // maximum number of instructions to execute (0 = unlimited)
	steps    int // instruction counter
}

// New creates a new VM ready to execute the given program.
func New(program *types.Program, io types.IO) *VM {
	return &VM{
		program:  program,
		vars:     make(map[string]types.Value),
		io:       io,
		maxSteps: DefaultMaxSteps,
	}
}

// SetMaxSteps sets the maximum number of instructions the VM will execute.
// Set to 0 for unlimited (not recommended).
func (vm *VM) SetMaxSteps(n int) {
	vm.maxSteps = n
}

// Run executes the program until it halts or encounters an error.
func (vm *VM) Run() error {
	for !vm.halted && vm.pc < len(vm.program.Instructions) {
		if vm.maxSteps > 0 {
			vm.steps++
			if vm.steps > vm.maxSteps {
				return types.LineError(vm.program.Instructions[vm.pc].Line,
					"execution limit exceeded (%d steps)", vm.maxSteps)
			}
		}

		inst := vm.program.Instructions[vm.pc]
		vm.pc++

		if err := vm.execute(inst); err != nil {
			return err
		}
	}
	return nil
}

// resolve converts an argument string to a Value.
// If the string matches a variable name, the variable's value is returned.
// Otherwise, it's parsed as a literal.
func (vm *VM) resolve(arg string) types.Value {
	// check if it's a variable
	if v, ok := vm.vars[arg]; ok {
		return v
	}

	// try to parse as literal
	return parseLiteral(arg)
}

// resolveVar resolves an argument that must be a defined variable.
// Used for arithmetic operands where using an undefined variable is an error.
func (vm *VM) resolveVar(arg string, line int) (types.Value, error) {
	v, ok := vm.vars[arg]
	if !ok {
		return types.Value{}, types.LineError(line, "undefined variable: %s", arg)
	}
	return v, nil
}

// getVar retrieves a variable's value, returning an error if undefined.
func (vm *VM) getVar(name string, line int) (types.Value, error) {
	v, ok := vm.vars[name]
	if !ok {
		return types.Value{}, types.LineError(line, "undefined variable: %s", name)
	}
	return v, nil
}

// resolveLabel resolves a label name or line number to an instruction index.
func (vm *VM) resolveLabel(target string, line int) (int, error) {
	// try as label name first
	if idx, ok := vm.program.Labels[target]; ok {
		return idx, nil
	}
	return 0, types.LineError(line, "undefined label: %s", target)
}

func parseLiteral(s string) types.Value {
	if s == "true" {
		return types.NewBool(true)
	}
	if s == "false" {
		return types.NewBool(false)
	}

	// try int
	var n int64
	neg := false
	i := 0
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		neg = s[0] == '-'
		i = 1
	}
	if i < len(s) {
		allDigits := true
		for j := i; j < len(s); j++ {
			if s[j] < '0' || s[j] > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && i < len(s) {
			for j := i; j < len(s); j++ {
				n = n*10 + int64(s[j]-'0')
			}
			if neg {
				n = -n
			}
			return types.NewInt(n)
		}
	}

	// try float
	hasDot := false
	for _, c := range s {
		if c == '.' {
			hasDot = true
			break
		}
	}
	if hasDot {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return types.NewFloat(f)
		}
	}

	// default to string
	return types.NewString(s)
}
