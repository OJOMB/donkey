package ast

import "github.com/OJOMB/donkey/internal/tokens"

// ExpressionIdentifier represents an identifier in the AST. It contains a token and a value for the identifier name.
type ExpressionIdentifier struct {
	// Token is the token associated with the identifier, which is typically a token.TypeIdent token.
	Token tokens.Token
	// Value is the name of the identifier, which is the lexeme of the token.
	Value string
}

func (i *ExpressionIdentifier) expressionNode()     {}
func (i *ExpressionIdentifier) TokenLexeme() string { return i.Token.Lexeme }

func (i *ExpressionIdentifier) String() string {
	return i.Value
}
