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
					var myFunction = fn(x) { return x + 1; };

					{ var z = 10;}
				`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
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
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
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
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "__foobar__"},
							Value: "__foobar__",
						},
						Value: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "FALSE", Lexeme: "false"},
							Value: false,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "myFunction"},
							Value: "myFunction",
						},
						Value: &ast.ExpressionLiteralFunction{
							Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
									Value: "x",
								},
							},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
											Value: 1,
										},
									},
								},
							}},
						},
					},
					&ast.StatementBlock{
						Statements: []ast.Statement{
							&ast.StatementBind{
								Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
								Name: &ast.ExpressionIdentifier{
									Token: tokens.Token{Type: "IDENT", Lexeme: "z"},
									Value: "z",
								},
								Value: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: "INT", Lexeme: "10"},
									Value: 10,
								},
							},
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
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralString{
							Token: tokens.Token{Type: "STRING", Lexeme: "kool"},
							Value: "kool",
						},
					},
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "TRUE", Lexeme: "true"},
							Value: true,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test rebind statements - no errors",
			input: `
				var x = 5;
				x = 10;
				x;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementRebind{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
							Value: 10,
						},
					},
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
							Value: "x",
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test while loop",
			input: `
				var result = 0;
				var i = 0;
				while (i < 10) {
					if (i % 2 == 0) {
						continue;
					}

					result = result + i;
					i = i + 1;

					if (result > 10) {
						break;
					}
				}`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "result"},
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "i"},
							Value: "i",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementWhile{
						Token: tokens.Token{Type: tokens.TypeWhile, Lexeme: "while"},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
							Operator: "<",
							Left: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
								Value: "i",
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
								Value: 10,
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementExpression{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.Token{Type: tokens.TypeEQ, Lexeme: "=="},
													Operator: "==",
													Left: &ast.ExpressionInfix{
														Token:    tokens.Token{Type: tokens.TypePercent, Lexeme: "%"},
														Operator: "%",
														Left: &ast.ExpressionIdentifier{
															Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
															Value: "i",
														},
														Right: &ast.ExpressionLiteralInteger{
															Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
															Value: 2,
														},
													},
													Right: &ast.ExpressionLiteralInteger{
														Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
														Value: 0,
													},
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
																Keyword: "continue",
															},
														},
													},
												},
											},
										},
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
											Value: "result",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
										Value: "i",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
										Right: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
											Value: 1,
										},
									},
								},
								&ast.StatementExpression{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.Token{Type: tokens.TypeGT, Lexeme: ">"},
													Operator: ">",
													Left: &ast.ExpressionIdentifier{
														Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
														Value: "result",
													},
													Right: &ast.ExpressionLiteralInteger{
														Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
														Value: 10,
													},
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.Token{Type: tokens.TypeBreak, Lexeme: "break"},
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.Token{Type: tokens.TypeBreak, Lexeme: "break"},
																Keyword: "break",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test increment and decrement operators are desugared correctly",
			input: `
				var x = 0;
				var y = 10;
				x++;
				y--;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
							Value: "y",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
							Value: 10,
						},
					},
					&ast.StatementRebind{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
								Value: "x",
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
								Value: 1,
							},
						},
					},
					&ast.StatementRebind{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
							Value: "y",
						},
						Value: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeMinus, Lexeme: "-"},
							Operator: "-",
							Left: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
								Value: "y",
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
								Value: 1,
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test function binding statement - no errors",
			input: `fn add(x, y) { return x + y; }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementFunctionBind{
						Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "add"},
							Value: "add",
						},
						Value: &ast.ExpressionLiteralFunction{
							Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
									Value: "x",
								},
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
									Value: "y",
								},
							},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
								},
							},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test builtin function len - no errors",
			input: `len("hello");`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "len"},
						Expression: &ast.ExpressionCall{
							Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
							Function: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "len"},
								Value: "len",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralString{
									Token: tokens.Token{Type: tokens.TypeString, Lexeme: "hello"},
									Value: "hello",
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test list literals - no errors",
			input: `
				var myList = [1, 2, 3];
				var myStrList = ["foo", "bar", "baz"];
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myList"},
							Value: "myList",
						},
						Value: &ast.ExpressionLiteralList{
							Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
							Elements: []ast.Expression{
								&ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
									Value: 1,
								},
								&ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
									Value: 2,
								},
								&ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
									Value: 3,
								},
							},
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myStrList"},
							Value: "myStrList",
						},
						Value: &ast.ExpressionLiteralList{
							Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
							Elements: []ast.Expression{
								&ast.ExpressionLiteralString{
									Token: tokens.Token{Type: tokens.TypeString, Lexeme: "foo"},
									Value: "foo",
								},
								&ast.ExpressionLiteralString{
									Token: tokens.Token{Type: tokens.TypeString, Lexeme: "bar"},
									Value: "bar",
								},
								&ast.ExpressionLiteralString{
									Token: tokens.Token{Type: tokens.TypeString, Lexeme: "baz"},
									Value: "baz",
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
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p, err := New(lexer.New(input, nil), nil)
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

	p, err := New(lexer.New(input, nil), nil)
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

	p, err := New(lexer.New(input, nil), nil)
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

		p, err := New(lexer.New(input, nil), nil)
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

		p, err := New(lexer.New(input, nil), nil)
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

		p, err := New(lexer.New(input, nil), nil)
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

		p, err := New(lexer.New(input, nil), nil)
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
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "122"},
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
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionInfix{
									Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
									Operator: "/",
									Left: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Value: 5,
									},
									Right: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
										Value: 10,
									},
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
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
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
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
						Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
							Operator: "*",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
								Operator: "+",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
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
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

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
						Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
										Operator: "<",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
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
		{
			name:  "if expression with elif and else - no errors",
			input: `if (x < y) { x } elif (x > y) { y } else { 5 }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
										Operator: "<",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
													Value: "x",
												},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeGT, Lexeme: ">"},
										Operator: ">",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
													Value: "y",
												},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Expression: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
											Value: 5,
										},
									},
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
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestFunctionLiterals(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "function literal with no parameters - no errors",
			input: `fn() { return 5; }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
						Expression: &ast.ExpressionLiteralFunction{
							Token:      tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Value: 5,
									},
								},
							}},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "function literal with parameters - no errors",
			input: `fn(x, y) { return x + y; }`,
			// this program is possible but functionally makes no sense - it would onbly make sense as the RHS of a StatementBind/StatementRebind.
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
						Expression: &ast.ExpressionLiteralFunction{
							Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
									Value: "x",
								},
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
									Value: "y",
								},
							},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
								},
							}},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParsingFunctionCalls(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "function call with args - no errors",
			input: `myFunction(2+3, param2, false);`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myFunction"},
						Expression: &ast.ExpressionCall{
							Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
							Function: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myFunction"},
								Value: "myFunction",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionInfix{
									Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
									Operator: "+",
									Left: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
										Value: 2,
									},
									Right: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
										Value: 3,
									},
								},
								&ast.ExpressionIdentifier{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "param2"},
									Value: "param2",
								},
								&ast.ExpressionLiteralBoolean{
									Token: tokens.Token{Type: tokens.TypeFalse, Lexeme: "false"},
									Value: false,
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
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseWhileLoop(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name: "test while loop",
			input: `
				var result = 0;
				var i = 0;
				while (i < 10) {
					result = result + i;
					i = i + 1;
				}`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "result"},
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "i"},
							Value: "i",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementWhile{
						Token: tokens.Token{Type: tokens.TypeWhile, Lexeme: "while"},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
							Operator: "<",
							Left: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
								Value: "i",
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
								Value: 10,
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
											Value: "result",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
										Value: "i",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
										Right: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
											Value: 1,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test while loop without condition",
			input: `
				var result = 0;
				var i = 0;
				while {
					if (i % 2 == 0) {
						continue;
					}

					result = result + i;
					i = i + 1;

					if (result > 10) {
						break;
					}
				}`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "result"},
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "i"},
							Value: "i",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
							Value: 0,
						},
					},
					&ast.StatementWhile{
						Token:     tokens.Token{Type: tokens.TypeWhile, Lexeme: "while"},
						Condition: nil,
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementExpression{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.Token{Type: tokens.TypeEQ, Lexeme: "=="},
													Operator: "==",
													Left: &ast.ExpressionInfix{
														Token:    tokens.Token{Type: tokens.TypePercent, Lexeme: "%"},
														Operator: "%",
														Left: &ast.ExpressionIdentifier{
															Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
															Value: "i",
														},
														Right: &ast.ExpressionLiteralInteger{
															Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
															Value: 2,
														},
													},
													Right: &ast.ExpressionLiteralInteger{
														Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
														Value: 0,
													},
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
																Keyword: "continue",
															},
														},
													},
												},
											},
										},
										Alternative: nil,
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
											Value: "result",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
										Value: "i",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
										Right: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
											Value: 1,
										},
									},
								},
								&ast.StatementExpression{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.Token{Type: tokens.TypeGT, Lexeme: ">"},
													Operator: ">",
													Left: &ast.ExpressionIdentifier{
														Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
														Value: "result",
													},
													Right: &ast.ExpressionLiteralInteger{
														Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
														Value: 10,
													},
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.Token{Type: tokens.TypeBreak, Lexeme: "break"},
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.Token{Type: tokens.TypeBreak, Lexeme: "break"},
																Keyword: "break",
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
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseForLoop(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name: "test for loop",
			input: `
				for (var i = 0; i < 10; i = i + 1) {
					if (i % 2 == 0) {
						continue;
					}

					result = result + i;
				}`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementFor{
						Token: tokens.Token{Type: tokens.TypeFor, Lexeme: "for"},
						Initializer: &ast.StatementBind{
							Token: tokens.Token{Type: tokens.TypeBind, Lexeme: "var"},
							Name: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: "IDENT", Lexeme: "i"},
								Value: "i",
							},
							Value: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
								Value: 0,
							},
						},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
							Operator: "<",
							Left: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
								Value: "i",
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
								Value: 10,
							},
						},
						Step: &ast.StatementRebind{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
							Name: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
								Value: "i",
							},
							Value: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
								Operator: "+",
								Left: &ast.ExpressionIdentifier{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
									Value: "i",
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
									Value: 1,
								},
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementExpression{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.Token{Type: tokens.TypeEQ, Lexeme: "=="},
													Operator: "==",
													Left: &ast.ExpressionInfix{
														Token:    tokens.Token{Type: tokens.TypePercent, Lexeme: "%"},
														Operator: "%",
														Left: &ast.ExpressionIdentifier{
															Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
															Value: "i",
														},
														Right: &ast.ExpressionLiteralInteger{
															Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
															Value: 2,
														},
													},
													Right: &ast.ExpressionLiteralInteger{
														Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
														Value: 0,
													},
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.Token{Type: tokens.TypeContinue, Lexeme: "continue"},
																Keyword: "continue",
															},
														},
													},
												},
											},
										},
										Alternative: nil,
									},
								},
								&ast.StatementRebind{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
									Name: &ast.ExpressionIdentifier{
										Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token: tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "result"},
											Value: "result",
										},
										Operator: "+",
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "i"},
											Value: "i",
										},
									},
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
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseExpressionLiteralList(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput ast.Expression
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "test empty list",
			input: `[]`,
			expectedOutput: &ast.ExpressionLiteralList{
				Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test list with multiple elements",
			input: `[1, 2 + 3, myFunction()]`,
			expectedOutput: &ast.ExpressionLiteralList{
				Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
				Elements: []ast.Expression{
					&ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
						Value: 1,
					},
					&ast.ExpressionInfix{
						Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
						Operator: "+",
						Left: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
							Value: 2,
						},
						Right: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
							Value: 3,
						},
					},
					&ast.ExpressionCall{
						Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
						Function: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myFunction"},
							Value: "myFunction",
						},
						Arguments: []ast.Expression{},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			// we know the program will be a single statement expression whose expression is the list literal, so we can directly access it here for ease of testing
			stmtExpr, ok := program.Statements[0].(*ast.StatementExpression)
			require.True(t, ok)

			assert.Equal(t, tc.expectedOutput, stmtExpr.Expression)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseExpressionListIndexing(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput ast.Statement
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "test simple list indexing",
			input: `myList[0]`,
			expectedOutput: &ast.StatementExpression{
				Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myList"},
				Expression: &ast.ExpressionIndex{
					Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
					Left: &ast.ExpressionIdentifier{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myList"},
						Value: "myList",
					},
					Index: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
						Value: 0,
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test list indexing with infix expression as index",
			input: `myList[1 + 2]`,
			expectedOutput: &ast.StatementExpression{
				Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myList"},
				Expression: &ast.ExpressionIndex{
					Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
					Left: &ast.ExpressionIdentifier{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myList"},
						Value: "myList",
					},
					Index: &ast.ExpressionInfix{
						Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
						Operator: "+",
						Left: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
							Value: 1,
						},
						Right: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
							Value: 2,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test indexing into a list literal",
			input: `[1, 2, 3][0]`,
			expectedOutput: &ast.StatementExpression{
				Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
				Expression: &ast.ExpressionIndex{
					Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
					Left: &ast.ExpressionLiteralList{
						Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
						Elements: []ast.Expression{
							&ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
								Value: 1,
							},
							&ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
								Value: 2,
							},
							&ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
								Value: 3,
							},
						},
					},
					Index: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
						Value: 0,
					},
				},
			},
			expectedErrs: []string{},
		},
		// {
		// 	// TODO: figure out a way to do this
		// 	name:  "test indexing into a nested list literal",
		// 	input: `[[1, 2, 3],[4,5,6]][0][1]`,
		// 	expectedOutput: &ast.StatementExpression{
		// 		Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 		Expression: &ast.ExpressionIndex{
		// 			Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 			Index: &ast.ExpressionLiteralInteger{
		// 				Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
		// 				Value: 1,
		// 			},
		// 			Left: &ast.ExpressionIndex{
		// 				Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 				Index: &ast.ExpressionLiteralInteger{
		// 					Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "0"},
		// 					Value: 0,
		// 				},
		// 				Left: &ast.ExpressionLiteralList{
		// 					Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 					Elements: []ast.Expression{
		// 						&ast.ExpressionLiteralList{
		// 							Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 							Elements: []ast.Expression{
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
		// 									Value: 1,
		// 								},
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
		// 									Value: 2,
		// 								},
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
		// 									Value: 3,
		// 								},
		// 							},
		// 						},
		// 						&ast.ExpressionLiteralList{
		// 							Token: tokens.Token{Type: tokens.TypeLBracket, Lexeme: "["},
		// 							Elements: []ast.Expression{
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
		// 									Value: 4,
		// 								},
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
		// 									Value: 5,
		// 								},
		// 								&ast.ExpressionLiteralInteger{
		// 									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "6"},
		// 									Value: 6,
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expectedErrs: []string{},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			// we know the program will be a single statement expression whose expression is the index expression, so we can directly access it here for ease of testing
			stmtExpr, ok := program.Statements[0].(*ast.StatementExpression)
			require.True(t, ok)

			assert.Equal(t, tc.expectedOutput, stmtExpr)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseBitwiseOps(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput ast.Statement
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "test bitwise shift left",
			input: `5 << 1`,
			expectedOutput: &ast.StatementExpression{
				Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
				Expression: &ast.ExpressionInfix{
					Token:    tokens.Token{Type: tokens.TypeBitwiseShiftLeft, Lexeme: "<<"},
					Operator: "<<",
					Left: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Value: 5,
					},
					Right: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
						Value: 1,
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "test bitwise shift right",
			input: `10 >> 2`,
			expectedOutput: &ast.StatementExpression{
				Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
				Expression: &ast.ExpressionInfix{
					Token:    tokens.Token{Type: tokens.TypeBitwiseShiftRight, Lexeme: ">>"},
					Operator: ">>",
					Left: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
						Value: 10,
					},
					Right: &ast.ExpressionLiteralInteger{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
						Value: 2,
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			// we know the program will be a single statement expression whose expression is the infix expression, so we can directly access it here for ease of testing
			stmtExpr, ok := program.Statements[0].(*ast.StatementExpression)
			require.True(t, ok)

			assert.Equal(t, tc.expectedOutput, stmtExpr)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParserParseExpressionLiteralMap(t *testing.T) {}
