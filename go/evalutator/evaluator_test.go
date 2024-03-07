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

func Test_ReturnStatement(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2*5; 9;", 10},
		{"return true;", true},
		{"return false;", false},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)

		switch tt.expected.(type) {
		case int:
			integer, _ := tt.expected.(int)
			testIntegerObject(t, evaluated, int64(integer))
		case bool:
			boolean, _ := tt.expected.(bool)
			testBooleanObject(t, evaluated, bool(boolean))
		}
	}
}

func Test_ErrorHandling(t *testing.T) {
	testData := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true+false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true+false; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}
			return 1;
		}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"hello" - "world"`, "unknown operator: STRING - STRING"},
		{`{"name": "monkey"}[fn(x) { x }]`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got %T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected %q, got %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func Test_LetStatement(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{`let a = "foobar"; a;`, "foobar"},
	}

	for _, tt := range testData {
		switch tt.expected.(type) {
		case int:
			i, _ := tt.expected.(int)
			testIntegerObject(t, testEval(tt.input), int64(i))
		case string:
			s, _ := tt.expected.(string)
			testStringObject(t, testEval(tt.input), s)
		}
	}
}

func Test_FunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function, got %T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("unexpected number of parameters. Parameters: %+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got %q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got %q", expectedBody, fn.Body.String())
	}
}

func Test_FunctionApplication(t *testing.T) {
	testData := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { return x*2; }; double(5);", 10},
		{"let add = fn(a, b) { a + b; }; add(5, 5)", 10},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range testData {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func Test_Closures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};
	let addTwo = newAdder(2);
	addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func Test_StringConcatenation(t *testing.T) {
	input := `"foo" + "_" + "bar"`
	expect := `foo_bar`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got %T (%+v)", evaluated, evaluated)
	}
	if str.Value != expect {
		t.Errorf("String has wrong value. got %q.", str.Value)
	}
}

func Test_BuiltInFunctions(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments to `len`. got=2, want=1"},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got %T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. \nexpected\t%q\nactual\t\t%q.", expected, errObj.Message)
			}
		}
	}
}

func Test_TypeArgument(t *testing.T) {
	testData := []struct {
		input    string
		expected string
	}{
		{`type("foo")`, "STRING"},
		{`type(false)`, "BOOLEAN"},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)
		tp, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("object is not String. got %T (%+v)", evaluated, evaluated)
			continue
		}
		if string(tp.Value) != tt.expected {
			t.Errorf("wrong type. \nexpected\t%q\nactual\t\t%q.", tt.expected, tp.Value)
		}
	}
}

func Test_ArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got %d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func Test_ArrayIndexExpressions(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1,2,3][1+1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
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

func Test_ArrayBuiltins(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{`len([1,2,3])`, 3},
		{`first([1,2,3])`, 1},
		{`last([1,2,3])`, 3},
		{`rest([1,2,3])`, "[2, 3]"},
		{`push([1,2,3], 4)`, "[1, 2, 3, 4]"},
	}

	for _, tt := range testData {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("response is not of type Array, got %T", evaluated)
				continue
			}
			if arr.Inspect() != tt.expected {
				t.Errorf("array mismatch. \nexpected\t%q\nactual\t\t%q", tt.expected, arr.Inspect())
			}
		}
	}
}

func Test_HashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10-9,
		two: 1+1,
		"thr"+"ee": 6/2,
		4:4,
		true: 5,
		false: 6
	}
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func Test_HashIndexExpressions(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
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
	return Eval(p.ParseProgram(), object.CreateEnvironment())
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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String, got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got %q, want %q", result.Value, expected)
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
