package ast

import "github.com/OJOMB/donkey/internal/tokens"

// ExpressionWhile represents a while loop expression in the AST. It contains a token, a condition expression, and a body statement block.
type ExpressionWhile struct {
	Token     tokens.Token
	Condition Expression
	Body      *StatementBlock
}

func (ew *ExpressionWhile) expressionNode()     {}
func (ew *ExpressionWhile) TokenLexeme() string { return "while" }

func (ew *ExpressionWhile) String() string {
	return "while " + ew.Condition.String() + " " + ew.Body.String()
}
