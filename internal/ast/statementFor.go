package ast

import (
	"github.com/OJOMB/donkey/internal/tokens"
)

// StatementFor represents a for loop statement in the AST.
// for (initializer; condition; step) { body }
type StatementFor struct {
	Token       tokens.Token    // The 'for' token
	Initializer Statement       // The initializer statement (e.g., var i = 0)
	Step        Statement       // The step statement (e.g., i = i + 1)
	Condition   Expression      // The condition expression
	Body        *StatementBlock // The body of the loop
}

func (s *StatementFor) statementNode() {}

func (s *StatementFor) TokenLexeme() string {
	return s.Token.Lexeme
}

func (s *StatementFor) String() string {
	result := s.TokenLexeme() + " ("
	if s.Initializer != nil {
		result += s.Initializer.String()
	}

	result += "; "
	if s.Condition != nil {
		result += s.Condition.String()
	}

	result += "; "
	if s.Step != nil {
		result += s.Step.String()
	}

	result += ") "
	if s.Body == nil {
		result += "{}"
		return result
	}

	return result + s.Body.String()
}

func (s *StatementFor) Type() NodeType {
	return NodeTypeStatementFor
}
