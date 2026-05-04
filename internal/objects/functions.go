package objects

import (
	"strings"

	"github.com/OJOMB/donkey/internal/ast"
)

type Function struct {
	Parameters []*ast.ExpressionIdentifier
	Body       *ast.StatementBlock
	Env        *Environment
}

func (f *Function) Type() Type {
	return TypeFunction
}

func (f *Function) Inspect() string {
	var b strings.Builder
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}

	b.WriteString("fn(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") {\n")
	b.WriteString(f.Body.String())
	b.WriteString("\n}")

	return b.String()
}
