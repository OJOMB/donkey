package vm

import "github.com/OJOMB/donkey/internal/bytecode"

type VMError struct {
	Message string
	OpCode  bytecode.Opcode
	IP      int
}

func (e *VMError) Error() string {
	return e.Message
}
