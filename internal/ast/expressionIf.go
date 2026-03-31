package ast

import (
	"strings"

	"github.com/OJOMB/monkey/internal/tokens"
)

type ExpressionIf struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *StatementBlock
	Alternative *StatementBlock
}

func (ei *ExpressionIf) expressionNode()     {}
func (ei *ExpressionIf) TokenLexeme() string { return ei.Token.Lexeme }

func (ei *ExpressionIf) String() string {
	// out := "if" + ei.Condition.String() + " " + ei.Consequence.String()
	var out = strings.Builder{}
	if _, err := out.WriteString("if"); err != nil {
		return "failed to write if expression string representation"
	}

	if _, err := out.WriteString(ei.Condition.String()); err != nil {
		return "failed to write if expression string representation"
	}

	if _, err := out.WriteString(" "); err != nil {
		return "failed to write if expression string representation"
	}

	if _, err := out.WriteString(ei.Consequence.String()); err != nil {
		return "failed to write if expression string representation"
	}

	if ei.Alternative != nil {
		if _, err := out.WriteString("else "); err != nil {
			return "failed to write if expression string representation"
		}
		if _, err := out.WriteString(ei.Alternative.String()); err != nil {
			return "failed to write if expression string representation"
		}
	}

	return out.String()
}
