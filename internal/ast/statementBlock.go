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
		if _, err := out.WriteString(stmt.String()); err != nil {
			return "failed to write statement block string representation"
		}
	}

	return out.String()
}
