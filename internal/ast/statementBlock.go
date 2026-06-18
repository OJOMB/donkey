package ast

import "strings"

type StatementBlock struct {
	Statements []Statement
}

func (sb *StatementBlock) statementNode()      {}
func (sb *StatementBlock) TokenLexeme() string { return "" }

func (sb *StatementBlock) String() string {
	var out = strings.Builder{}
	for _, stmt := range sb.Statements {
		_, _ = out.WriteString(stmt.String())
	}

	return out.String()
}

func (sb *StatementBlock) Type() NodeType {
	return NodeTypeStatementBlock
}
