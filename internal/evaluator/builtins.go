package evaluator

import (
	"strings"

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
		},
		Implementation: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return newError("len expects exactly one argument, got %d", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.String:
				return &objects.Integer{Value: int(len([]rune(arg.Value)))}
			case *objects.List:
				return &objects.Integer{Value: int(len(arg.Elements))}
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
		},
		Implementation: func(args ...objects.Object) objects.Object {
			strBuilder := strings.Builder{}
			for i, arg := range args {
				strBuilder.WriteString(arg.Inspect())
				if i < len(args)-1 {
					strBuilder.WriteString(" ")
				}
			}

			return &objects.String{Value: strBuilder.String()}
		},
		Name: "print",
	},
	// cdr returns a new list containing all elements of the input list except the first one.
	// If the input list is empty, cdr returns an empty list. If the input is not a list, cdr returns an error.
	"cdr": {
		Fn: objects.Function{
			Parameters: []*ast.ExpressionIdentifier{
				{Value: "arg"},
			},
			Body: &ast.StatementBlock{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Expression: &ast.ExpressionCall{
							Function: &ast.ExpressionIdentifier{Value: "cdr"},
							Arguments: []ast.Expression{
								&ast.ExpressionIdentifier{Value: "arg"},
							},
						},
					},
				},
			},
		},
		Implementation: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return newError("cdr expects exactly one argument, got %d", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.List:
				if len(arg.Elements) == 0 {
					return &objects.List{Elements: []objects.Object{}}
				}

				return &objects.List{Elements: arg.Elements[1:]}
			default:
				return newError("cdr not supported for type %s", args[0].Type())
			}
		},
		Name: "cdr",
	},
	// car returns the first element of a list.
	// If the input list is empty, car returns nowt. If the input is not a list, car returns an error.
	"car": {
		Fn: objects.Function{
			Parameters: []*ast.ExpressionIdentifier{
				{Value: "arg"},
			},
			Body: &ast.StatementBlock{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Expression: &ast.ExpressionCall{
							Function: &ast.ExpressionIdentifier{Value: "car"},
							Arguments: []ast.Expression{
								&ast.ExpressionIdentifier{Value: "arg"},
							},
						},
					},
				},
			},
		},
		Implementation: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return newError("car expects exactly one argument, got %d", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.List:
				if len(arg.Elements) == 0 {
					return &objects.Nowt{}
				}

				return arg.Elements[0]
			default:
				return newError("car not supported for type %s", args[0].Type())
			}
		},
		Name: "car",
	},
	// cons takes an element and a list and returns a new list with the element added to the front of the list.
	// If the second argument is not a list, cons returns an error.
	"cons": {
		Fn: objects.Function{
			Parameters: []*ast.ExpressionIdentifier{
				{Value: "head"},
				{Value: "tail"},
			},
			Body: &ast.StatementBlock{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Expression: &ast.ExpressionCall{
							Function: &ast.ExpressionIdentifier{Value: "cons"},
							Arguments: []ast.Expression{
								&ast.ExpressionIdentifier{Value: "head"},
								&ast.ExpressionIdentifier{Value: "tail"},
							},
						},
					},
				},
			},
		},
		Implementation: func(args ...objects.Object) objects.Object {
			if len(args) != 2 {
				return newError("cons expects exactly two arguments, got %d", len(args))
			}

			switch tail := args[1].(type) {
			case *objects.List:
				return &objects.List{Elements: append([]objects.Object{args[0]}, tail.Elements...)}
			default:
				return newError("cons expects second argument to be a list, got %s", args[1].Type())
			}
		},
		Name: "cons",
	},
}
