package ast

import (
	"strings"

	"github.com/OJOMB/donkey/internal/tokens"
)

type ConditionalBranch struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *StatementBlock
}

type ExpressionIf struct {
	Branches    []ConditionalBranch // first = if, rest = elif
	Alternative *StatementBlock     // optional else
}

func (ei *ExpressionIf) expressionNode()     {}
func (ei *ExpressionIf) TokenLexeme() string { return "if" }

func (ei *ExpressionIf) String() string {
	// out := "if" + ei.Condition.String() + " " + ei.Consequence.String()
	var out = strings.Builder{}
	if _, err := out.WriteString("if"); err != nil {
		return "failed to write if expression string representation"
	}

	for i, branch := range ei.Branches {
		if i > 0 {
			if _, err := out.WriteString("elif"); err != nil {
				return "failed to write if expression string representation"
			}
		}

		if _, err := out.WriteString(branch.Condition.String()); err != nil {
			return "failed to write if expression string representation"
		}
		if _, err := out.WriteString(" "); err != nil {
			return "failed to write if expression string representation"
		}
		if _, err := out.WriteString(branch.Consequence.String()); err != nil {
			return "failed to write if expression string representation"
		}
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
