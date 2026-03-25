package ast

import (
	"fmt"

	"github.com/OJOMB/monkey/internal/tokens"
)

// stringFmtLetStatement is the format string used for representing a let statement in the AST when converting it to a string.
// let <identifier> = <expression>;
const stringFmtLetStatement = "let %s = %s;"

// StatementLet represents a let statement in the AST. It contains a token, an identifier for the variable name, and an expression for the value.
type StatementLet struct {
	// Token is the token associated with the let statement, which is typically a token.TokenTypeLet token.
	Token tokens.Token
	// Name is the identifier for the variable name being declared in the let statement.
	// LHS of the let statement, which is the variable name being declared.
	Name *ExpressionIdentifier
	// Value is the expression that represents the value being assigned to the variable in the let statement.
	// RHS of the let statement, which is the expression representing the value being assigned to the variable.
	Value Expression
}

func (ls *StatementLet) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the let statement.
func (ls *StatementLet) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *StatementLet) String() string {
	return fmt.Sprintf(stringFmtLetStatement, ls.Name.String(), ls.Value.String())
}
