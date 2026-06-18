package ast

import "github.com/OJOMB/donkey/internal/tokens"

// StatementWhile represents a while loop statement in the AST. It contains a token, a condition expression, and a body statement block.
type StatementWhile struct {
	Token     tokens.Token
	Condition Expression
	Body      *StatementBlock
}

// statementNode is a marker method to indicate that StatementWhile is a statement node in the AST.
func (ew *StatementWhile) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the while statement, which is typically "while".
func (ew *StatementWhile) TokenLexeme() string { return "while" }

// Type returns the NodeType of the StatementWhile node, which is NodeTypeStatementWhile.
func (ew *StatementWhile) String() string {
	return "while " + ew.Condition.String() + " " + ew.Body.String()
}

func (ew *StatementWhile) Type() NodeType {
	return NodeTypeStatementWhile
}
