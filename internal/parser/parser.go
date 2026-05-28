package parser

import (
	"fmt"

	"github.com/Moonhalf383/pigsh_kernel/internal/lexer"
	"github.com/Moonhalf383/pigsh_kernel/internal/types"
)

// Parse converts a token stream into a Program (instruction list + label table).
// Each line of tokens is expected to start with a keyword token.
func Parse(tokens []lexer.Token) (*types.Program, error) {
	prog := &types.Program{
		Labels: make(map[string]int),
	}

	i := 0
	for i < len(tokens) {
		tok := &tokens[i]
		if tok.Kind != lexer.TokenKeyword {
			return nil, fmt.Errorf("Error at line %d: expected instruction, got %q", tok.Line, tok.Literal)
		}

		op := tok.Literal
		line := tok.Line
		i++

		// collect remaining tokens on this line (same line number)
		var args []string
		for i < len(tokens) && tokens[i].Line == line {
			args = append(args, tokens[i].Literal)
			i++
		}

		// handle label definition: "label <name>"
		if op == "label" {
			if len(args) < 1 {
				return nil, types.LineError(line, "label requires a name")
			}
			name := args[0]
			if _, exists := prog.Labels[name]; exists {
				return nil, types.LineError(line, "duplicate label: %s", name)
			}
			prog.Labels[name] = len(prog.Instructions)
			// still emit a nop instruction so line indices stay consistent
			prog.Instructions = append(prog.Instructions, types.Instruction{
				Line: line,
				Op:   "label",
				Args: args,
			})
			continue
		}

		prog.Instructions = append(prog.Instructions, types.Instruction{
			Line: line,
			Op:   op,
			Args: args,
		})
	}

	return prog, nil
}
