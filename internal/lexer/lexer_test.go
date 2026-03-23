package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OJOMB/monkey/internal/tokens"
)

func TestNextToken1(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput []tokens.Token
	}

	var testCases = []testCase{
		{
			name: "test next token with all token types",
			input: `
				let five = 5;
				let ten10 = 10;

				let add = fn(x, y) {
					let sum = x + y;

					return sum;
				};

				if (five < 10) {
					let __one = false;
				} elif (ten10 > 10) {
				 	let __two = true;
				} elif (ten10 == 10) {
					let __three = false;
				} elif (ten10 != 10) {
					let __four = true;
				} else {
					let __five = false;
				}

				let result = add(five, ten);`,
			expectedOutput: []tokens.Token{
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "five"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeInt, Lexeme: "5"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten10"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "add"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeFunction, Lexeme: "fn"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "x"},
				{Type: tokens.TokenTypeComma, Lexeme: ","},
				{Type: tokens.TokenTypeIdent, Lexeme: "y"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "sum"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeIdent, Lexeme: "x"},
				{Type: tokens.TokenTypePlus, Lexeme: "+"},
				{Type: tokens.TokenTypeIdent, Lexeme: "y"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeReturn, Lexeme: "return"},
				{Type: tokens.TokenTypeIdent, Lexeme: "sum"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeIf, Lexeme: "if"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "five"},
				{Type: tokens.TokenTypeLT, Lexeme: "<"},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "__one"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeFalse, Lexeme: "false"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeElif, Lexeme: "elif"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten10"},
				{Type: tokens.TokenTypeGT, Lexeme: ">"},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "__two"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeTrue, Lexeme: "true"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeElif, Lexeme: "elif"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten10"},
				{Type: tokens.TokenTypeEquality, Lexeme: "=="},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "__three"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeFalse, Lexeme: "false"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeElif, Lexeme: "elif"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten10"},
				{Type: tokens.TokenTypeNotEqual, Lexeme: "!="},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "__four"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeTrue, Lexeme: "true"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeElse, Lexeme: "else"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "__five"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeFalse, Lexeme: "false"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "result"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeIdent, Lexeme: "add"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "five"},
				{Type: tokens.TokenTypeComma, Lexeme: ","},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeEOF, Lexeme: ""},
			},
		},
	}

	for i, tc := range testCases {
		lex := New(tc.input)

		// call NextToken until we get an EOF token - assuming the lexer is working correctly we should get the expected output tokens in order
		var toks []tokens.Token
		for {
			tok := lex.NextToken()
			toks = append(toks, tok)

			if tok.Type == tokens.TokenTypeEOF {
				break
			}
		}

		assert.Equal(t, tc.expectedOutput, toks, "test case %d failed", i)
	}
}
