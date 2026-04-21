package evaluator

import (
	"fmt"
	"strings"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/tokens"
	"github.com/OJOMB/donkey/pkg/logs"
)

var (
	// Nowt is the singleton Nowt object that represents the absence of a value in the Donkey programming language.
	Nowt = &objects.Nowt{}
	// True is the singleton Boolean object that represents the boolean value true in the Donkey programming language.
	True = &objects.Boolean{Value: true}
	// False is the singleton Boolean object that represents the boolean value false in the Donkey programming language.
	False = &objects.Boolean{Value: false}
)

type Evaluator struct {
	logger logs.Logger
}

func New(l logs.Logger) *Evaluator {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Evaluator{logger: l.With("component", "evaluator")}
}

// Eval evaluates the given AST node and returns the resulting object.
func (e *Evaluator) Eval(node ast.Node) objects.Object {
	switch nt := node.(type) {
	case *ast.ExpressionLiteralInteger, *ast.ExpressionLiteralBoolean, *ast.ExpressionLiteralString, *ast.ExpressionLiteralFunction:
		return e.evalLiteral(nt)
	case *ast.Program:
		return e.evalStatements(nt)
	case *ast.StatementExpression:
		return e.Eval(nt.Expression)
	case *ast.ExpressionPrefix:
		right := e.Eval(nt.Right)
		if right == nil {
			e.logger.Error("prefix operator right-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return Nowt
		}

		switch nt.Token.Type {
		case tokens.TypeBang:
			return e.evalExpressionPrefixBang(right)
		case tokens.TypeMinus:
			return e.evalExpressionPrefixMinus(right)
		default:
			e.logger.Error("unsupported prefix operator", "operator", nt.Token.Lexeme)
			return Nowt
		}
	case *ast.ExpressionInfix:
		l := e.Eval(nt.Left)
		if l == nil {
			e.logger.Error("infix operator left-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return Nowt
		}

		r := e.Eval(nt.Right)
		if r == nil {
			e.logger.Error("infix operator right-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return Nowt
		}

		return e.evalExpressionInfix(nt.Operator, l, r)
	default:
		e.logger.Warn("unsupported AST node type", "type", fmt.Sprintf("%T", nt))
		return Nowt
	}
}

func (e *Evaluator) evalStatements(program *ast.Program) objects.Object {
	var result objects.Object
	for i, stmt := range program.Statements {
		e.logger.Debug("evaluating statement", "index", i, "statement", stmt.String())
		result = e.Eval(stmt)
	}

	return result
}

func (e *Evaluator) evalExpressionPrefixBang(right objects.Object) objects.Object {
	if right.Type() != objects.TypeBoolean {
		e.logger.Warn("unsupported operand type for ! operator", "type", right.Type())
		return Nowt
	}

	switch right {
	case True:
		return False
	case False:
		return True
	default:
		return Nowt
	}
}

func (e *Evaluator) evalExpressionPrefixMinus(right objects.Object) objects.Object {
	if right.Type() != objects.TypeInteger {
		e.logger.Warn("unsupported operand type for - operator", "type", right.Type())
		return Nowt
	}

	value := right.(*objects.Integer).Value
	return &objects.Integer{Value: -value}
}

func (e *Evaluator) evalLiteral(node ast.Node) objects.Object {
	switch nt := node.(type) {
	case *ast.ExpressionLiteralInteger:
		return &objects.Integer{Value: nt.Value}
	case *ast.ExpressionLiteralBoolean:
		if nt.Value {
			return True
		}
		return False
	case *ast.ExpressionLiteralString:
		return &objects.String{Value: nt.Value}
	case *ast.ExpressionLiteralFunction:
		// function literals are not evaluated to a value until they are called, so we return a Nowt object for now
		return Nowt
	default:
		e.logger.Warn("unsupported literal type", "type", fmt.Sprintf("%T", nt))
		return Nowt
	}
}

func (e *Evaluator) evalExpressionInfix(operator string, left, right objects.Object) objects.Object {
	switch {
	case left.Type() == objects.TypeInteger && right.Type() == objects.TypeInteger:
		return e.evalExpressionInfixInteger(operator, left.(*objects.Integer), right.(*objects.Integer))
	case left.Type() == objects.TypeBoolean && right.Type() == objects.TypeBoolean:
		return e.evalExpressionInfixBoolean(operator, left.(*objects.Boolean), right.(*objects.Boolean))
	case left.Type() == objects.TypeString && right.Type() == objects.TypeString:
		return e.evalExpressionInfixString(operator, left.(*objects.String), right.(*objects.String))
	default:
		e.logger.Warn("unsupported operand types for infix operator", "operator", operator, "leftType", left.Type(), "rightType", right.Type())
		return Nowt
	}
}

func (e *Evaluator) evalExpressionInfixInteger(operator string, left, right *objects.Integer) objects.Object {
	switch operator {
	case "+":
		return &objects.Integer{Value: left.Value + right.Value}
	case "-":
		return &objects.Integer{Value: left.Value - right.Value}
	case "*":
		return &objects.Integer{Value: left.Value * right.Value}
	case "/":
		if right.Value == 0 {
			e.logger.Warn("division by zero")
			return Nowt
		}

		return &objects.Integer{Value: left.Value / right.Value}
	case "%":
		if right.Value == 0 {
			e.logger.Warn("modulo by zero")
			return Nowt
		}

		return &objects.Integer{Value: left.Value % right.Value}
	case "^":
		if right.Value == 0 {
			return &objects.Integer{Value: 1}
		}

		result := 1
		for i := 0; i < right.Value; i++ {
			result *= left.Value
		}

		return &objects.Integer{Value: result}
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	case "<":
		return &objects.Boolean{Value: left.Value < right.Value}
	case ">":
		return &objects.Boolean{Value: left.Value > right.Value}
	case "<=":
		return &objects.Boolean{Value: left.Value <= right.Value}
	case ">=":
		return &objects.Boolean{Value: left.Value >= right.Value}
	case "&":
		return &objects.Integer{Value: left.Value & right.Value}
	case "|":
		return &objects.Integer{Value: left.Value | right.Value}
	default:
		e.logger.Warn("unsupported infix operator for integers", "operator", operator)
		return Nowt
	}
}

func (e *Evaluator) evalExpressionInfixBoolean(operator string, left, right *objects.Boolean) objects.Object {
	switch operator {
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	case "&&":
		return &objects.Boolean{Value: left.Value && right.Value}
	case "||":
		return &objects.Boolean{Value: left.Value || right.Value}
	default:
		e.logger.Warn("unsupported infix operator for booleans", "operator", operator)
		return Nowt
	}
}

func (e *Evaluator) evalExpressionInfixString(operator string, left, right *objects.String) objects.Object {
	switch operator {
	case "+":
		return &objects.String{Value: left.Value + right.Value}
	case "-":
		// TODO: not overly convinced about this one
		return &objects.String{Value: strings.TrimSuffix(left.Value, right.Value)}
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	default:
		e.logger.Warn("unsupported infix operator for strings", "operator", operator)
		return Nowt
	}
}
