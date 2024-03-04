package evalutator

import (
	"testing"

	"alde.nu/mint/lexer"
	"alde.nu/mint/object"
	"alde.nu/mint/parser"
)

func Test_EvalIntegerExpression(t *testing.T) {
	testData := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

/// Helper functions /////////////////////////////////////////////////

func testEval(input string) object.Object {
	l := lexer.Create(input)
	p := parser.Create(l)
	return Eval(p.ParseProgram())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %d, want %d", result.Value, expected)
		return false
	}
	return true
}
