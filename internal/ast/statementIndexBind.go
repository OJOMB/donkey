package ast

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/tokens"
)

// stringFmtIndexBindStatement is the format string used for representing an index bind statement in the AST when converting it to a string.
// <ExpressionIndex> = <Expression>;
const stringFmtIndexBindStatement = "%s = %s;"

// StatementIndexBind represents an index bind statement in the AST.
// It contains a token, an index expression for the key, and an expression for the value.
// For example, in the index bind statement `m["key"] = "new value";`, the token would be the "=" token, the Key would be an ExpressionIndex representing the key, and the Value would be an Expression representing the new value.
type StatementIndexBind struct {
	// Token is the "=" token in the index bind statement.
	Token tokens.Token
	// Left is any expression that resolves to an indexable
	Left *ExpressionIndex
	// Right is the expression that represents the value being assigned to the key in the index bind statement.
	// RHS of the index bind statement, which is the expression representing the value being assigned to the key.
	Right Expression
}

func (ls *StatementIndexBind) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the index bind statement.
func (ls *StatementIndexBind) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *StatementIndexBind) String() string {
	return fmt.Sprintf(stringFmtIndexBindStatement, ls.Left.String(), ls.Right.String())
}
