package compiler

import (
	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/bytecode"
	"github.com/OJOMB/donkey/internal/objects"
)

type Bytecode struct {
	Instructions []bytecode.Instruction
	Constants    []objects.Object
}

type Compiler struct {
	instructions []bytecode.Instruction
	constants    []objects.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: make([]bytecode.Instruction, 0),
		constants:    make([]objects.Object, 0),
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, s := range n.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.StatementExpression:
		if err := c.Compile(n.Expression); err != nil {
			return err
		}
	case *ast.ExpressionInfix:
		if err := c.Compile(n.Left); err != nil {
			return err
		}
		if err := c.Compile(n.Right); err != nil {
			return err
		}
	case *ast.ExpressionLiteralInteger:
		i := &objects.Integer{Value: n.Value}
		c.emit(bytecode.OpCodeConstant, c.addConstant(i))
	}

	return nil
}

// addConstant adds a new object constant to the internal constant store and returns its index position.
// The index position can be used as the constants unique identifier.
func (c *Compiler) addConstant(obj objects.Object) int {
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(ins bytecode.Instruction) int {
	c.instructions = append(c.instructions, ins)

	return len(c.instructions) - 1
}

// emit takes a bytecode opcode and its operands, creates a new instruction, and appends it to the compiler's instruction list.
// It returns the index of the newly added instruction in the instruction list.
func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	ins := bytecode.NewInstruction(op, operands...)
	pos := c.addInstruction(ins)

	return pos
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
