package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/tokens"
)

func TestParseStatements(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name: "test bind statements - no errors",
			input: `
					var x = 5;
					var y = "hello";
					var __foobar__ = false;
				`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TokenTypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TokenTypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "y"},
							Value: "y",
						},
						Value: &ast.ExpressionLiteralString{
							Token: tokens.Token{Type: "STRING", Lexeme: "hello"},
							Value: "hello",
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TokenTypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "__foobar__"},
							Value: "__foobar__",
						},
						Value: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "FALSE", Lexeme: "false"},
							Value: false,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test return statements",
			input: `
				return 5;
				return "kool";
				return true;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionLiteralString{
							Token: tokens.Token{Type: "STRING", Lexeme: "kool"},
							Value: "kool",
						},
					},
					&ast.ReturnStatement{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						ReturnValue: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "TRUE", Lexeme: "true"},
							Value: true,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "assignment statements",
			input: `
				x = 5;
				y = "kool";
				__foobar__ = true;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{},
			},
			expectedErrs: []string{
				"unexpected token: =",
				"unexpected token: =",
				"unexpected token: =",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
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

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
	require.True(t, ok)

	intLiteral, ok := stmt.Expression.(*ast.ExpressionLiteralInteger)
	require.True(t, ok)
	assert.Equal(t, 5, intLiteral.Value)
	assert.Equal(t, "5", intLiteral.TokenLexeme())
	assert.Equal(t, "5", intLiteral.String())
}

func TestStringExpression(t *testing.T) {
	input := `"hello";`

	p, err := New(lexer.New(input), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
	require.True(t, ok)

	strLiteral, ok := stmt.Expression.(*ast.ExpressionLiteralString)
	require.True(t, ok)
	assert.Equal(t, "hello", strLiteral.Value)
	assert.Equal(t, "hello", strLiteral.TokenLexeme())
	assert.Equal(t, "hello", strLiteral.String())
}

func TestExpressionPrefix(t *testing.T) {
	t.Run("prefix expression: !", func(t *testing.T) {
		input := "!5;"

		p, err := New(lexer.New(input), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionPrefix{}, stmt.Expression)

		prefixExp, ok := stmt.Expression.(*ast.ExpressionPrefix)
		require.True(t, ok)
		assert.Equal(t, "!", prefixExp.Operator)
		assert.Equal(t, 5, prefixExp.Right.(*ast.ExpressionLiteralInteger).Value)
		assert.Equal(t, "!", prefixExp.TokenLexeme())
	})

	t.Run("prefix expression: -", func(t *testing.T) {
		input := "-15;"

		p, err := New(lexer.New(input), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionPrefix{}, stmt.Expression)

		prefixExp, ok := stmt.Expression.(*ast.ExpressionPrefix)
		require.True(t, ok)
		assert.Equal(t, "-", prefixExp.Operator)
		assert.Equal(t, 15, prefixExp.Right.(*ast.ExpressionLiteralInteger).Value)
		assert.Equal(t, "-", prefixExp.TokenLexeme())
	})
}

func TestExpressionStatementBool(t *testing.T) {
	t.Run("simple boolean literal: true", func(t *testing.T) {
		input := "true;"

		p, err := New(lexer.New(input), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionLiteralBoolean{}, stmt.Expression)

		boolExp, ok := stmt.Expression.(*ast.ExpressionLiteralBoolean)
		require.True(t, ok)
		assert.Equal(t, true, boolExp.Value)
		assert.Equal(t, "true", boolExp.TokenLexeme())
	})

	t.Run("simple boolean literal: false", func(t *testing.T) {
		input := "false;"

		p, err := New(lexer.New(input), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionLiteralBoolean{}, stmt.Expression)

		boolExp, ok := stmt.Expression.(*ast.ExpressionLiteralBoolean)
		require.True(t, ok)
		assert.Equal(t, false, boolExp.Value)
		assert.Equal(t, "false", boolExp.TokenLexeme())
	})
}

func TestParsingInfixExpressions(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "simple infix expressions - no errors",
			input: `5 + 122;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TokenTypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "122"},
								Value: 122,
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "slightly more complex infix expression - no errors",
			input: `5 + 5 / 10 * 4;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TokenTypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TokenTypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionInfix{
									Token:    tokens.Token{Type: tokens.TokenTypeForwardSlash, Lexeme: "/"},
									Operator: "/",
									Left: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
										Value: 5,
									},
									Right: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "10"},
										Value: 10,
									},
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "slightly more complex infix expression - no errors",
			input: `5 * 5 + 10 / 4;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TokenTypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TokenTypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TokenTypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "grouped infix expression - no errors",
			input: `(5 + 5) * (10 / 4);`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TokenTypeLParen, Lexeme: "("},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TokenTypeAsterisk, Lexeme: "*"},
							Operator: "*",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TokenTypePlus, Lexeme: "+"},
								Operator: "+",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TokenTypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TokenTypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestIfExpression(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "if expression without else - no errors",
			input: `if (x < y) { x }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TokenTypeIf, Lexeme: "if"},
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TokenTypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TokenTypeLT, Lexeme: "<"},
										Operator: "<",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TokenTypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TokenTypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TokenTypeIdent, Lexeme: "x"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TokenTypeIdent, Lexeme: "x"},
													Value: "x",
												},
											},
										},
									},
								},
							},
							Alternative: nil,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}
