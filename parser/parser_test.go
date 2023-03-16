package parser

import (
	"BubblePL/ast"
	"BubblePL/lexer"
	"fmt"
	"log"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x  5;
let = 10;
let 838 383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseError(t, p)
	if program == nil {
		t.Fatalf("Parse Program Error\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("get wrong statements, expected=%d, got=%d", 3, len(program.Statements))
	}
	tests := []struct {
		expectedNameStr string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for idx, test := range tests {
		testLetStatement(t, program.Statements[idx], test.expectedNameStr)
	}
}

func testLetStatement(t *testing.T, st ast.Statement, name string) {
	if st.ToLiteral() != "let" {
		t.Fatalf("Statement.ToLiteral error, expected=%s, got=%s", "let", st.ToLiteral())
	}

	letSt, ok := st.(*ast.LetStatement)
	if !ok {
		t.Fatalf("got wrong type, expected=*ast.LetStatement, got=%T", letSt)
	}

	if letSt.Name.ToLiteral() != name {
		t.Fatalf("get wrong name, expected='%s', got='%s'", name, letSt.Name.ToLiteral())
	}

}

func checkParseError(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}
	for _, errMsg := range errs {
		t.Errorf(errMsg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993 233;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStatement.ToLiteral() != "return" {
			t.Errorf("returnStatement.ToLiteral not 'return'. got=%s", returnStatement.ToLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements got wrong length got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Errorf("Statement.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("Identifier.Value is not 'foobar', got=%s", ident.Value)
	}

	if ident.ToLiteral() != "foobar" {
		t.Errorf("Identifier.ToLiteral is not 'foobar', got=%s", ident.Value)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements got wrong length got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Statement.Expression is not ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if ident.Value != 5 {
		t.Errorf("IntegerLiteral.Value is not 5, got=%d", ident.Value)
	}

	if ident.ToLiteral() != "5" {
		t.Errorf("Identifier.ToLiteral is not '5', got=%d", ident.Value)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	exp, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Experssion is not ast.IntegerLiteral. got=%T", il)
		return false
	}

	if exp.Value != value {
		t.Errorf("ast.IntegerLiteral.Value is not %d. got=%d", value, exp.Value)
		return false
	}
	if exp.ToLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("ast.IntegerLiteral.ToLiteral is not '%d'. got=%s", value, exp.ToLiteral())
		return false
	}
	return true
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements got wrong length got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s' *ast.PrefixExpression. got='%s'", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParseError(t, parser)
		if len(program.Statements) != 1 {
			log.Fatalf("program.Statement got wrong length. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.InfixExpression. got=%T", program.Statements[0])
		}
		if !testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParseError(t, parser)
		programText := program.String()
		if programText != tt.expected {
			t.Errorf("wrong program text expected=%s, got=%s", tt.expected, programText)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression is not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("Identifier.Value is not '%s'. got=%s", value, ident.Value)
		return false
	}
	if ident.ToLiteral() != value {
		t.Errorf("Identifier.ToLiteral() is not '%s'. got=%s", value, ident.ToLiteral())
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
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp is not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expression is not *ast.InfixExpression. got=%T", exp)
		return false
	}
	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}
	if infixExp.Operator != operator {
		t.Errorf("infixExp.Operator is not '%s'. got=%s", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("stmt.Expression is not *ast.Boolean. got=%T", exp)
		return false
	}
	if boExp.Value != value {
		t.Errorf("ast.Boolean is not %t. got=%t", value, boExp.Value)
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) {x}`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParseError(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements got wrong length got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence.Statements got wrong length got=%d", len(exp.Consequence.Statements))
	}
	stmt, ok = exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative should be nil")
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}
