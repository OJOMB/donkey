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
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
