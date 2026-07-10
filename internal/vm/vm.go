package vm

import (
	"github.com/OJOMB/donkey/internal/bytecode"
	"github.com/OJOMB/donkey/internal/compiler"
	"github.com/OJOMB/donkey/internal/objects"
)

const (
	stackSize = 2048
)

type VM struct {
	instructions []bytecode.Instruction
	constants    []objects.Object
	stack        []any
	sp           int // stack pointer``
}

func New(bc *compiler.Bytecode) *VM {
	return &VM{
		instructions: bc.Instructions,
		constants:    bc.Constants,
		stack:        make([]any, stackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for ip, instruction := range vm.instructions {
		switch bytecode.Opcode(instruction[0]) {
		case bytecode.OpCodeConstant:
			constIndex := int(instruction[1])<<8 | int(instruction[2])
			constant := vm.constants[constIndex]
			vm.push(constant)
		default:
			return &VMError{Message: "unknown opcode: %d", OpCode: bytecode.Opcode(instruction[0]), IP: ip}
		}
	}

	return nil
}

func (vm *VM) StackTop() any {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) push(obj any) {
	vm.stack[vm.sp] = obj
	vm.sp++
}

func (vm *VM) pop() any {
	vm.sp--
	obj := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil

	return obj
}
