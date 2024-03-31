package vm

import (
	"fmt"
	"testing"

	"alde.nu/mint/ast"
	"alde.nu/mint/compiler"
	"alde.nu/mint/lexer"
	"alde.nu/mint/object"
	"alde.nu/mint/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func Test_IntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 3},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		if err := vm.Run(); err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		if err := testIntegerObject(int64(expected), actual); err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.Create(input)
	p := parser.Create(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T  (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value.\n\tgot=%d\n\twant=%d", result.Value, expected)
	}

	return nil
}
