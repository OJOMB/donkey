package evaluator

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
)

type builtinLib map[string]*objects.BuiltinFunction

var builtins = builtinLib{
	"len": {
		Fn: objects.Function{
			Parameters: []*ast.ExpressionIdentifier{
				{Value: "arg"},
			},
			Body: &ast.StatementBlock{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Expression: &ast.ExpressionCall{
							Function: &ast.ExpressionIdentifier{Value: "len"},
							Arguments: []ast.Expression{
								&ast.ExpressionIdentifier{Value: "arg"},
							},
						},
					},
				},
			},
			Env: nil, // Built-in functions don't have an environment
		},
		Implementation: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return newError("len expects exactly one argument, got %d", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.String:
				return &objects.Integer{Value: int(len(arg.Value))}
			// case *objects.Array:
			// 	return &objects.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("len not supported for type %s", args[0].Type())
			}
		},
		Name: "len",
	},
	"print": {
		Fn: objects.Function{
			Parameters: []*ast.ExpressionIdentifier{
				{Value: "arg"},
			},
			Body: &ast.StatementBlock{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Expression: &ast.ExpressionCall{
							Function: &ast.ExpressionIdentifier{Value: "print"},
							Arguments: []ast.Expression{
								&ast.ExpressionIdentifier{Value: "arg"},
							},
						},
					},
				},
			},
			Env: nil, // Built-in functions don't have an environment
		},
		Implementation: func(args ...objects.Object) objects.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return &objects.Nowt{}
		},
		Name: "print",
	},
}
