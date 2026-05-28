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
// Int and Float convert directly; String attempts numeric parsing; Bool returns 0.
func (v Value) AsNumber() float64 {
	switch v.Type {
	case TypeInt:
		return float64(v.Int)
	case TypeFloat:
		return v.Float
	case TypeString:
		s := v.String
		if s == "" {
			return 0
		}
		i := 0
		neg := false
		if s[0] == '-' || s[0] == '+' {
			neg = s[0] == '-'
			i = 1
		}
		hasDigit := false
		hasDot := false
		var intPart int64
		for ; i < len(s); i++ {
			c := s[i]
			if c >= '0' && c <= '9' {
				hasDigit = true
				intPart = intPart*10 + int64(c-'0')
			} else if c == '.' && !hasDot {
				hasDot = true
				// parse the fractional part and return as float
				i++
				var frac float64
				place := 0.1
				for ; i < len(s); i++ {
					if s[i] < '0' || s[i] > '9' {
						return 0
					}
					frac += float64(s[i]-'0') * place
					place *= 0.1
				}
				result := float64(intPart) + frac
				if neg {
					result = -result
				}
				return result
			} else {
				return 0
			}
		}
		if !hasDigit {
			return 0
		}
		if neg {
			return float64(-intPart)
		}
		return float64(intPart)
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
	Line   int      // original source line number (1-based)
	Op     string   // operation code: "var", "add", "print", etc.
	Args   []string // arguments as raw strings
	Quoted []bool   // per-arg: true if originally a quoted string literal
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
