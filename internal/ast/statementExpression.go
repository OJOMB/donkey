package ast

import "github.com/OJOMB/monkey/internal/tokens"

// StatementExpression represents a statement in the AST that is just an expression. It contains a token and an expression.
// For example, in the Monkey programming language, an expression statement could be something like "5 + 5;" where the expression is "5 + 5" or "foobar;" where the expression is "foobar".
type StatementExpression struct {
	// Token is the token associated with the expression statement, which is typically the first token of the expression.
	Token tokens.Token
	// Expression is the expression contained within the expression statement.
	Expression Expression
}

func (es *StatementExpression) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the expression statement.
func (es *StatementExpression) TokenLexeme() string { return es.Token.Lexeme }

func (es *StatementExpression) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}
