package ast

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/tokens"
)

// stringFmtLetStatement is the format string used for representing a let statement in the AST when converting it to a string.
// var <ExpressionIdentifier> = <Expression>;
const stringFmtLetStatement = "var %s = %s;"

// StatementBind represents a var statement in the AST. It contains a token, an identifier for the variable name, and an expression for the value.
// For example, in the var statement "var x = 5;", the token would be the "var" token, the Name would be an ExpressionIdentifier representing "x", and the Value would be an ExpressionLiteralInteger representing "5".
// Or we might have var x = y; in which case the Value would be an ExpressionIdentifier representing "y".
type StatementBind struct {
	// Token is the token associated with the var statement, which is typically a token.TypeVar token.
	Token tokens.Token
	// Name is the identifier for the variable name being declared in the var statement.
	// LHS of the var statement, which is the variable name being declared.
	Name *ExpressionIdentifier
	// Value is the expression that represents the value being assigned to the variable in the var statement.
	// RHS of the var statement, which is the expression representing the value being assigned to the variable.
	Value Expression
}

func (ls *StatementBind) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the var statement.
func (ls *StatementBind) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *StatementBind) String() string {
	return fmt.Sprintf(stringFmtLetStatement, ls.Name.String(), ls.Value.String())
}
