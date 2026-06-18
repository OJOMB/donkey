package ast

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/tokens"
)

// stringFmtRebindStatement is the format string used for representing a rebind statement in the AST when converting it to a string.
// <ExpressionIdentifier> = <Expression>;
const stringFmtRebindStatement = "%s = %s;"

// StatementRebind represents a rebind statement in the AST. It contains a token, an identifier for the variable name, and an expression for the value.
// For example, in the rebind statement "x = 5;", the token would be the "=" token, the Name would be an ExpressionIdentifier representing "x", and the Value would be an ExpressionLiteralInteger representing "5".
// Or we might have x = y; in which case the Value would be an ExpressionIdentifier representing "y".
type StatementRebind struct {
	// Token is the token
	Token tokens.Token
	// Name is the identifier for the variable name being rebound in the rebind statement.
	// LHS of the rebind statement, which is the variable name being rebound.
	Name *ExpressionIdentifier
	// Value is the expression that represents the value being assigned to the variable in the rebind statement.
	// RHS of the rebind statement, which is the expression representing the value being assigned to the variable.
	Value Expression
}

func (ls *StatementRebind) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the rebind statement.
func (ls *StatementRebind) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *StatementRebind) String() string {
	return fmt.Sprintf(stringFmtRebindStatement, ls.Name.String(), ls.Value.String())
}

func (ls *StatementRebind) Type() NodeType {
	return NodeTypeStatementRebind
}
