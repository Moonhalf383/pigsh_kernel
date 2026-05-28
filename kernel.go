package pigsh

import (
	"github.com/Moonhalf383/pigsh_kernel/internal/lexer"
	"github.com/Moonhalf383/pigsh_kernel/internal/parser"
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
	"github.com/Moonhalf383/pigsh_kernel/internal/vm"
)

// IO defines the interface for output and input during script execution.
// Callers implement this to wire pigsh output to their own systems (QQ bot, CLI, etc.).
type IO = types.IO

// Options configures VM behavior.
type Options struct {
	// MaxSteps limits the number of instructions executed. 0 means unlimited.
	// Default is 1,000,000 if not set.
	MaxSteps int
}

// Run parses and executes pigsh source code with default options.
// Returns nil on success, or an error describing what went wrong (with line numbers).
func Run(source string, io IO) error {
	return RunWithOptions(source, io, Options{})
}

// RunWithOptions parses and executes pigsh source code with custom options.
func RunWithOptions(source string, io IO, opts Options) error {
	tokens := lexer.Tokenize(source)
	if len(tokens) == 0 {
		return nil
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	machine := vm.New(program, io)
	if opts.MaxSteps != 0 {
		machine.SetMaxSteps(opts.MaxSteps)
	}
	return machine.Run()
}
