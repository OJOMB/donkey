package vm

import "github.com/OJOMB/donkey/internal/compiler"

type VM struct{}

func New(bc *compiler.Bytecode) *VM {
	return &VM{}
}

func (vm *VM) Run() error {
	return nil
}

func (vm *VM) StackTop() any {
	return nil
}
