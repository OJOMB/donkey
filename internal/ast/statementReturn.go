package ast

import (
	"fmt"

	"github.com/OJOMB/monkey/internal/tokens"
)

const stringFmtReturnStatement = "return %s;"

// ReturnStatement represents a return statement in the AST. It contains a token and an expression for the return value.
type ReturnStatement struct {
	// Token is the token associated with the return statement, which is typically a token.TokenTypeReturn token.
	Token tokens.Token
	// ReturnValue is the expression that represents the value being returned by the return statement.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the return statement.
func (rs *ReturnStatement) TokenLexeme() string {
	return rs.Token.Lexeme
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf(stringFmtReturnStatement, rs.ReturnValue.String())
}
