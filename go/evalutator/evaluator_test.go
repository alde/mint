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
		{"-10", -10},
		{"-5", -5},
		{"5-10", -5},
		{"10+5", 15},
		{"10/5", 2},
		{"10*5", 50},
		{"5+5+5+5-10", 10},
		{"2*2*2*2*2", 32},
		{"-50 + 100 + -50", 0},
		{"5 + 2 * 10", 25},
		{"5 * 2 + 10", 20},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 - 10", 50},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func Test_EvalBooleanExpression(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{"false", false},
		{"true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) == false", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test_BangOperator(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{"!false", true},
		{"!true", false},
		{"!!false", false},
		{"!!true", true},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test_IfElseExpression(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1<2) { 10 }", 10},
		{"if (1>2) { 10 }", nil},
		{"if (1>2) { 10 } else { 20 }", 20},
		{"if (1<2) { 10 } else { 20 }", 10},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %t, want %t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("expected object to be NULL, got %T (%+v)", obj, obj)
		return false
	}
	return true
}
