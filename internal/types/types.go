package types

import "fmt"

// ValueType represents the type of a pigsh value.
type ValueType int

const (
	TypeInt ValueType = iota
	TypeFloat
	TypeString
	TypeBool
)

// Value is a tagged union representing all pigsh runtime values.
type Value struct {
	Type   ValueType
	Int    int64
	Float  float64
	String string
	Bool   bool
}

func NewInt(v int64) Value    { return Value{Type: TypeInt, Int: v} }
func NewFloat(v float64) Value { return Value{Type: TypeFloat, Float: v} }
func NewString(v string) Value { return Value{Type: TypeString, String: v} }
func NewBool(v bool) Value     { return Value{Type: TypeBool, Bool: v} }

// AsNumber converts the value to float64 for arithmetic.
// Int and Float convert directly; String and Bool return 0.
func (v Value) AsNumber() float64 {
	switch v.Type {
	case TypeInt:
		return float64(v.Int)
	case TypeFloat:
		return v.Float
	default:
		return 0
	}
}

// AsInt converts the value to int64.
func (v Value) AsInt() int64 {
	switch v.Type {
	case TypeInt:
		return v.Int
	case TypeFloat:
		return int64(v.Float)
	default:
		return 0
	}
}

// AsBool returns the boolean representation.
func (v Value) AsBool() bool {
	switch v.Type {
	case TypeBool:
		return v.Bool
	case TypeInt:
		return v.Int != 0
	case TypeFloat:
		return v.Float != 0
	case TypeString:
		return v.String != ""
	default:
		return false
	}
}

// Format returns the string representation for print output.
func (v Value) Format() string {
	switch v.Type {
	case TypeInt:
		return fmt.Sprintf("%d", v.Int)
	case TypeFloat:
		return fmt.Sprintf("%g", v.Float)
	case TypeString:
		return v.String
	case TypeBool:
		if v.Bool {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// IsNumeric returns true if the value is Int or Float.
func (v Value) IsNumeric() bool {
	return v.Type == TypeInt || v.Type == TypeFloat
}

// Instruction represents a single parsed pigsh instruction.
type Instruction struct {
	Line int      // original source line number (1-based)
	Op   string   // operation code: "var", "add", "print", etc.
	Args []string // arguments as raw strings
}

// Program is the result of parsing a pigsh source.
type Program struct {
	Instructions []Instruction
	Labels       map[string]int // label name → instruction index
}

// IO defines the interface for output and input during execution.
type IO interface {
	Print(value string)
	Input(prompt string) string
}

// LineError creates an error prefixed with the source line number.
func LineError(line int, format string, a ...any) error {
	return fmt.Errorf("Error at line %d: %s", line, fmt.Sprintf(format, a...))
}
