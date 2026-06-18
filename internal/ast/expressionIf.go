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
	var out = strings.Builder{}

	// NB WriteString never returns an error when writing to a strings.Builder, it's purely for interface compatibility with io.Writer.
	_, _ = out.WriteString("if")
	for i, branch := range ei.Branches {
		if i > 0 {
			_, _ = out.WriteString("elif")
		}

		_, _ = out.WriteString(branch.Condition.String())
		_, _ = out.WriteString(" ")
		_, _ = out.WriteString(branch.Consequence.String())
	}

	if ei.Alternative != nil {
		_, _ = out.WriteString("else ")
		_, _ = out.WriteString(ei.Alternative.String())
	}

	return out.String()
}

func (ei *ExpressionIf) Type() NodeType {
	return NodeTypeExpressionIf
}
