package pigsh_test

import (
	"strings"
	"testing"

	pigsh "github.com/Moonhalf383/pigsh_kernel"
)

// testIO captures print output and provides scripted input.
type testIO struct {
	output  []string
	inputs  []string
	inputIdx int
}

func (t *testIO) Print(value string) {
	t.output = append(t.output, value)
}

func (t *testIO) Input(prompt string) string {
	if t.inputIdx < len(t.inputs) {
		val := t.inputs[t.inputIdx]
		t.inputIdx++
		return val
	}
	return ""
}

func (t *testIO) Output() string {
	return strings.Join(t.output, "\n")
}

func TestHelloWorld(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`print hello`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "hello" {
		t.Errorf("expected 'hello', got %q", io.Output())
	}
}

func TestFibonacci(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var a 0
var b 1
var n 10
var i 0
label loop
print a
var tmp 0
mov tmp b
add b a b
mov a tmp
inc i
blt i n loop
halt`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "0\n1\n1\n2\n3\n5\n8\n13\n21\n34"
	if io.Output() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, io.Output())
	}
}

func TestSubroutine(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`call say_hi
call say_bye
halt
label say_hi
print hello
ret
label say_bye
print goodbye
ret`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "hello\ngoodbye"
	if io.Output() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, io.Output())
	}
}

func TestAccumulator(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var total 0
var i 1
var n 10
label loop
add total total i
inc i
ble i n loop
print total
halt`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "55" {
		t.Errorf("expected '55', got %q", io.Output())
	}
}

func TestInput(t *testing.T) {
	io := &testIO{inputs: []string{"pig"}}
	err := pigsh.Run(`input name who_are_you
print name`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "pig" {
		t.Errorf("expected 'pig', got %q", io.Output())
	}
}

func TestArithmetic(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var x 0
add x x 5
add x x 10
mul x x 2
print x`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "30" {
		t.Errorf("expected '30', got %q", io.Output())
	}
}

func TestStringConcat(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var greeting ""
add greeting "hello" "world"
print greeting`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", io.Output())
	}
}

func TestDivisionByZero(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var x 10
var z 0
div x x z`, io)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Errorf("expected 'division by zero' error, got: %v", err)
	}
}

func TestUndefinedVariable(t *testing.T) {
	io := &testIO{}
	// inc on an undeclared variable should error
	err := pigsh.Run(`inc x`, io)
	if err == nil {
		t.Fatal("expected undefined variable error")
	}
	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("expected 'undefined variable' error, got: %v", err)
	}
}

func TestPrintLiteral(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`print undefined_var`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "undefined_var" {
		t.Errorf("expected 'undefined_var', got %q", io.Output())
	}
}

func TestPrintVariable(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`var x 42
print x`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "42" {
		t.Errorf("expected '42', got %q", io.Output())
	}
}

func TestStackUnderflow(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`pop x`, io)
	if err == nil {
		t.Fatal("expected stack underflow error")
	}
	if !strings.Contains(err.Error(), "stack underflow") {
		t.Errorf("expected 'stack underflow' error, got: %v", err)
	}
}

func TestCallStackUnderflow(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`ret`, io)
	if err == nil {
		t.Fatal("expected call stack underflow error")
	}
	if !strings.Contains(err.Error(), "call stack underflow") {
		t.Errorf("expected 'call stack underflow' error, got: %v", err)
	}
}

func TestNestedCalls(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`call a
halt
label a
print in_a
call b
print back_a
ret
label b
print in_b
ret`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "in_a\nin_b\nback_a"
	if io.Output() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, io.Output())
	}
}

func TestComments(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`# this is a comment
print hello
# another comment
print world`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "hello\nworld"
	if io.Output() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, io.Output())
	}
}

func TestEmptySource(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(``, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(io.output) != 0 {
		t.Errorf("expected no output, got %q", io.Output())
	}
}

func TestHalt(t *testing.T) {
	io := &testIO{}
	err := pigsh.Run(`print before
halt
print after`, io)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if io.Output() != "before" {
		t.Errorf("expected 'before', got %q", io.Output())
	}
}
