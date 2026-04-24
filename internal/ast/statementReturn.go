package ast

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/tokens"
)

const stringFmtReturnStatement = "return %s;"

// StatementReturn represents a return statement in the AST. It contains a token and an expression for the return value.
type StatementReturn struct {
	// Token is the token associated with the return statement, which is typically a token.TypeReturn token.
	Token tokens.Token
	// Value is the expression that represents the value being returned by the return statement.
	Value Expression
}

func (rs *StatementReturn) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the return statement.
func (rs *StatementReturn) TokenLexeme() string {
	return rs.Token.Lexeme
}

func (rs *StatementReturn) String() string {
	return fmt.Sprintf(stringFmtReturnStatement, rs.Value.String())
}
