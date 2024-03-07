package parser

import (
	"fmt"
	"testing"

	"alde.nu/mint/ast"
	"alde.nu/mint/lexer"
	"alde.nu/mint/token"
)

func Test_LetStatements(t *testing.T) {
	testData := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func Test_ReturnStatements(t *testing.T) {
	testData := []struct {
		input    string
		expected interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return foobar;", "foobar"},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", returnStatement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got %q", returnStatement.TokenLiteral())
		}
		val := returnStatement.ReturnValue
		if !testLiteralExpression(t, val, tt.expected) {
			return
		}
	}
}

func Test_IdentifierExpression(t *testing.T) {
	input := "boo;"
	l := lexer.Create(input)
	p := Create(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have expected number of statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, "boo")
}

func Test_BooleanExpression(t *testing.T) {
	input := "true;"
	l := lexer.Create(input)
	p := Create(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have expected number of statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Boolean)

	if !ok {
		t.Fatalf("expression isn't *ast.Boolean. got=%T", stmt.Expression)
	}

	if ident.Value != true {
		t.Errorf("ident.Value isn't %v. got=%v", true, ident.Value)
	}

	if ident.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral() isn't %s. got=%s", "true", ident.TokenLiteral())
	}
}

func Test_IntegerExpression(t *testing.T) {
	input := "54;"
	l := lexer.Create(input)
	p := Create(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program doesn't have expected number of statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, 54)
}

func Test_ParsingPrefixExpressions(t *testing.T) {
	testData := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)

		exp := basicParsingChecks(t, program, 1, &ast.PrefixExpression{})

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func Test_ParsingInfixExpressions(t *testing.T) {
	testData := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range testData {
		l := lexer.Create(tt.input)
		p := Create(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program doesn't have expected number of statements, expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func Test_OperatorPrecedenceParsing(t *testing.T) {
	testData := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3>5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1+(2+3)+4", "((1 + (2 + 3)) + 4)"},
		{"(3+4)*5", "((3 + 4) * 5)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b*c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2*3, 4+5, add(6, 7 + 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 + 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func Test_IfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	program := initTests(t, input)

	exp := basicParsingChecks(t, program, 1, &ast.IfExpression{})

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func Test_IfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	program := initTests(t, input)

	exp := basicParsingChecks(t, program, 1, &ast.IfExpression{})

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was nil. got=%+v", exp.Alternative)
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statement. got=%d\n", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func Test_FunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y}`

	program := initTests(t, input)
	function := basicParsingChecks(t, program, 1, &ast.FunctionLiteral{})

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. Want 2, got %d\n", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements count is wrong. Expected %d, got %d\n", 1, len(function.Body.Statements))
	}

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statement[0] is not an ast.ExpressionStatement, got=%T\n", function.Body.Statements[0])
	}

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func Test_ParameterParsin(t *testing.T) {
	testData := []struct {
		input    string
		expected []string
	}{
		{input: `fn() {}`, expected: []string{}},
		{input: `fn(x) {}`, expected: []string{"x"}},
		{input: `fn(x, y, z) {}`, expected: []string{"x", "y", "z"}},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected) {
			t.Fatalf("wrong number of parameters, expected %d got %d", len(tt.expected), len(function.Parameters))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func Test_CallExpression(t *testing.T) {
	input := `add(1, 2+3, 4 * 5)`
	program := initTests(t, input)

	exp := basicParsingChecks(t, program, 1, &ast.CallExpression{})

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("expected 3 arguments, got %d (%s)", len(exp.Arguments), exp.Arguments)
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "+", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "*", 5)
}

func Test_CallArgumentParsing(t *testing.T) {
	testData := []struct {
		input    string
		expected []string
	}{
		{input: `add()`, expected: []string{}},
		{input: `add(x)`, expected: []string{"x"}},
		{input: `add(x, y, z)`, expected: []string{"x", "y", "z"}},
	}

	for _, tt := range testData {
		program := initTests(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		call := stmt.Expression.(*ast.CallExpression)

		if len(call.Arguments) != len(tt.expected) {
			t.Fatalf("wrong number of parameters, expected %d got %d", len(tt.expected), len(call.Arguments))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, call.Arguments[i], ident)
		}
	}
}

func Test_StringLiteralExpression(t *testing.T) {
	input := `"hello world"`
	program := initTests(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got %q.", "hello world", literal.Value)
	}
}

func Test_ParseArrayLiteral(t *testing.T) {
	input := `[1, 2*2, "foo", 3+3]`

	program := initTests(t, input)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp is not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 4 {
		t.Fatalf("len(array.Elements) not 4. got %d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testStringLiteral(t, array.Elements[2], "foo")
	testInfixExpression(t, array.Elements[3], 3, "+", 3)
}

func Test_ParsingIndexExpressions(t *testing.T) {
	input := `myArray[1+1]`
	program := initTests(t, input)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func Test_ParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	program := initTests(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Paris has wrong length, got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func Test_ParseEmptyHashLiteral(t *testing.T) {
	input := "{}"
	program := initTests(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Fatalf("hash.Paris has wrong length, got=%d", len(hash.Pairs))
	}
}

func Test_ParseHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	program := initTests(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Paris has wrong length, got=%d", len(hash.Pairs))
	}

	testData := map[string]func(ast.Expression){
		"one":   func(e ast.Expression) { testInfixExpression(t, e, 0, "+", 1) },
		"two":   func(e ast.Expression) { testInfixExpression(t, e, 10, "-", 8) },
		"three": func(e ast.Expression) { testInfixExpression(t, e, 15, "/", 5) },
	}

	for key, value := range hash.Pairs {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := testData[lit.String()]
		if !ok {
			t.Errorf("no test function for key %q found", lit.String())
			continue
		}

		testFunc(value)
	}
}

/// Helper functions /////////////////////////////////////////////////

func initTests(t *testing.T, input string) *ast.Program {
	l := lexer.Create(input)
	p := Create(l)
	defer checkParserErrors(t, p)
	return p.ParseProgram()
}

func basicParsingChecks[T ast.Expression](t *testing.T, program *ast.Program, expectedStatements int, kind T) T {
	if len(program.Statements) != expectedStatements {
		t.Fatalf("program.Statements does not contain %d statement. got=%d", expectedStatements, (program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(T)
	if !ok {
		t.Fatalf("stmt.Expression is not %T. got=%T", kind, stmt.Expression)
	}

	return exp
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not '%s'. got=%q", token.LET, s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name.TokenLiteral() not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d. got=%s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, il ast.Expression, value string) bool {
	str, ok := il.(*ast.StringLiteral)
	if !ok {
		t.Errorf("il not *ast.StringLiteral. got=%T", il)
		return false
	}

	if str.Value != value {
		t.Errorf("string.Value not %s. got=%s", value, str.Value)
		return false
	}

	if str.TokenLiteral() != value {
		t.Errorf("string.TokenLiteral() not %s. got=%s", value, str.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, il ast.Expression, value bool) bool {
	boolean, ok := il.(*ast.Boolean)
	if !ok {
		t.Errorf("il not *ast.Boolean. got=%T", il)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser found %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExpression.Left, left) {
		return false
	}

	if opExpression.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpression.Right, right) {
		return false
	}

	return true
}
