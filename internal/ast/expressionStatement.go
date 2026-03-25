package ast

import "github.com/OJOMB/monkey/internal/tokens"

// ExpressionStatement represents an expression statement in the AST. It contains a token and an expression.
type ExpressionStatement struct {
	// Token is the token associated with the expression statement, which is typically the first token of the expression.
	Token tokens.Token
	// Expression is the expression contained within the expression statement.
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the expression statement.
func (es *ExpressionStatement) TokenLexeme() string { return es.Token.Lexeme }

func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}
