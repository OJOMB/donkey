package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/monkey/internal/ast"
	"github.com/OJOMB/monkey/internal/lexer"
	"github.com/OJOMB/monkey/internal/tokens"
)

func TestLetStmts(t *testing.T) {
	type testCase struct {
		Name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			Name: "test let statements - no errors",
			input: `
				let x = 5;
				let y = 10;
				let __foobar__ = 838383;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementLet{
						Token: tokens.Token{Type: "LET", Lexeme: "let"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.LiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementLet{
						Token: tokens.Token{Type: "LET", Lexeme: "let"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "y"},
							Value: "y",
						},
					},
					&ast.StatementLet{
						Token: tokens.Token{Type: "LET", Lexeme: "let"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "__foobar__"},
							Value: "__foobar__",
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			Name: "test return statements",
			input: `
				return 5;
				return 10;
				return 993322;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: "5",
						},
					},
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "INT", Lexeme: "10"},
							Value: "10",
						},
					},
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "INT", Lexeme: "993322"},
							Value: "993322",
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p, err := New(lexer.New(input), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	ident, ok := stmt.Expression.(*ast.ExpressionIdentifier)
	require.True(t, ok)
	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLexeme())
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	p, err := New(lexer.New(input), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	intLiteral, ok := stmt.Expression.(*ast.LiteralInteger)
	require.True(t, ok)
	assert.Equal(t, 5, intLiteral.Value)
	assert.Equal(t, "5", intLiteral.TokenLexeme())
}

func TestExpressionPrefix(t *testing.T) {
	input := "!5;"

	p, err := New(lexer.New(input), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	prefixExp, ok := stmt.Expression.(*ast.ExpressionPrefix)
	require.True(t, ok)
	assert.Equal(t, "!", prefixExp.Operator)
	assert.Equal(t, 5, prefixExp.Right.(*ast.LiteralInteger).Value)
	assert.Equal(t, "!", prefixExp.TokenLexeme())
}
