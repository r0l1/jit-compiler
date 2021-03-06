package asm

/*
	Instructions

	The below instructions implement a simple Instruction interface that allows
	them to be combined, optimised and encoded into machine code.

	x86_64 is a bit interesting in that instructions with different sets of
	operands might require subtly different machine code opcodes, even though
	they do the same thing functionally speaking, so the below instructions
	make use of a thing called OpcodeMaps that will match the given arguments
	to an appropriate opcode. For more on OpcodeMaps see asm/opcodes/
*/

import (
	"fmt"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/asm/opcodes"
	"github.com/bspaans/jit/lib"
)

func ADD(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("add", opcodes.ADD, 2, dest, src)
}
func CALL(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("call", opcodes.CALL, 1, dest)
}
func CMP(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cmp", opcodes.CMP, 2, dest, src)
}

// Convert signed integer to scalar double-precision floating point (float64)
func CVTSI2SD(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cvtsi2sd", opcodes.CVTSI2SD, 2, dest, src)
}

// Convert double precision float to signed integer
func CVTTSD2SI(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cvttsd2si", opcodes.CVTTSD2SI, 2, dest, src)
}
func DEC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("dec", opcodes.DEC, 1, dest)
}
func DIV(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("div", opcodes.DIV, 2, dest, src)
}
func INC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("inc", opcodes.INC, 1, dest)
}
func JNE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("jne", opcodes.JNE, 1, dest)
}
func JMP(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("jmp", opcodes.JMP, 1, dest)
}
func LEA(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("lea", opcodes.LEA, 2, dest, src)
}
func MOV(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("mov", opcodes.MOV, 2, dest, src)
}
func MUL(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("mul", opcodes.MUL, 2, dest, src)
}
func POP(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("pop", opcodes.POP, 1, dest)
}
func PUSH(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("push", opcodes.PUSH, 1, dest)
}
func PUSHFQ() lib.Instruction {
	return opcodes.OpcodeToInstruction("pushfq", opcodes.PUSHFQ, 0)
}
func RETURN() lib.Instruction {
	return opcodes.OpcodeToInstruction("return", opcodes.RETURN, 0)
}
func SETA(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETA, 1, dest)
}
func SETAE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setae", opcodes.SETAE, 1, dest)
}
func SETB(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETB, 1, dest)
}
func SETBE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setae", opcodes.SETBE, 1, dest)
}
func SETC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETC, 1, dest)
}
func SETE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sete", opcodes.SETE, 1, dest)
}
func SETNE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setne", opcodes.SETNE, 1, dest)
}
func SUB(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sub", opcodes.SUB, 2, dest, src)
}
func SYSCALL() lib.Instruction {
	return opcodes.OpcodeToInstruction("syscall", opcodes.SYSCALL, 0)
}
func XOR(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("xor", opcodes.XOR, 2, dest, src)
}

func init() {
	b, err := lib.CompileInstruction([]lib.Instruction{
		MOV(encoding.Uint64(0), encoding.Rax),
		MOV(encoding.Uint64(0), encoding.Rcx),
		MOV(encoding.Uint64(0), encoding.Rdx),
		MOV(encoding.Uint64(0), encoding.Rbx),
		MOV(encoding.Uint64(0), encoding.Rbp),
		MOV(encoding.Uint64(0), encoding.Rsi),
		MOV(encoding.Uint64(0), encoding.Rdi),
		MOV(encoding.Uint32(0), &encoding.DisplacedRegister{encoding.Rsp, 8}),
		MOV(encoding.Uint64(0xffff), encoding.Rdi),
		INC(encoding.Rax),
		CMP(encoding.Rdi, encoding.Rax),
		JNE(encoding.Uint8(0xf9)),
		CMP(encoding.Rdi, encoding.Rax),
		SETE(encoding.Al),
		MOV(encoding.Uint64(123), encoding.Rcx),
		ADD(encoding.Rcx, encoding.Rax),
		ADD(encoding.Uint32(2), encoding.Rax),
		PUSH(encoding.Rax),
		POP(encoding.Rax),
		PUSHFQ(),
		POP(encoding.Rdx),
		MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
		RETURN(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
