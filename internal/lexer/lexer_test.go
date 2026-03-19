package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OJOMB/monkey/internal/lexer/tokens"
)

func TestNextToken(t *testing.T) {
	in := `
		let five = 5;
		let ten = 10;

		let add = fn(x, y) {
			x + y;
		};

		let result = add(five, ten);`

	expectedOut := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.TokenTypeLet, "let"},
		{tokens.TokenTypeIdent, "five"},
		{tokens.TokenTypeAssign, "="},
		{tokens.TokenTypeInt, "5"},
		{tokens.TokenTypeSemicolon, ";"},

		{tokens.TokenTypeLet, "let"},
		{tokens.TokenTypeIdent, "ten"},
		{tokens.TokenTypeAssign, "="},
		{tokens.TokenTypeInt, "10"},
		{tokens.TokenTypeSemicolon, ";"},

		{tokens.TokenTypeLet, "let"},
		{tokens.TokenTypeIdent, "add"},
		{tokens.TokenTypeAssign, "="},
		{tokens.TokenTypeFunction, "fn"},
		{tokens.TokenTypeLParen, "("},
		{tokens.TokenTypeIdent, "x"},
		{tokens.TokenTypeComma, ","},
		{tokens.TokenTypeIdent, "y"},
		{tokens.TokenTypeRParen, ")"},
		{tokens.TokenTypeLBrace, "{"},
		{tokens.TokenTypeIdent, "x"},
		{tokens.TokenTypePlus, "+"},
		{tokens.TokenTypeIdent, "y"},
		{tokens.TokenTypeSemicolon, ";"},
		{tokens.TokenTypeRBrace, "}"},
		{tokens.TokenTypeSemicolon, ";"},

		{tokens.TokenTypeLet, "let"},
		{tokens.TokenTypeIdent, "result"},
		{tokens.TokenTypeAssign, "="},
		{tokens.TokenTypeIdent, "add"},
		{tokens.TokenTypeLParen, "("},
		{tokens.TokenTypeIdent, "five"},
		{tokens.TokenTypeComma, ","},
		{tokens.TokenTypeIdent, "ten"},
		{tokens.TokenTypeRParen, ")"},
		{tokens.TokenTypeSemicolon, ";"},
	}

	lex := New(in)
	for i, tt := range expectedOut {
		tok := lex.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type, "tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Lexeme, "tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Lexeme)
	}
}
