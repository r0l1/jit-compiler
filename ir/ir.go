package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/expr"
	"github.com/bspaans/jit/ir/shared"
	. "github.com/bspaans/jit/ir/shared"
	. "github.com/bspaans/jit/ir/statements"
	"github.com/bspaans/jit/lib"
)

func Compile(stmts []IR) (lib.MachineCode, error) {
	ctx := NewIRContext()
	result := []uint8{}
	fmt.Println(".data:")
	for _, stmt := range stmts {
		currentOffset := ctx.DataSectionOffset + len(ctx.DataSection)
		if err := stmt.AddToDataSection(ctx); err != nil {
			return nil, err
		}
		if len(ctx.DataSection) != currentOffset-2 {
			fmt.Printf("0x%x-0x%x (0x%x): %s\n",
				currentOffset,
				len(ctx.DataSection)+ctx.DataSectionOffset,
				len(ctx.DataSection)+ctx.DataSectionOffset-currentOffset, stmt.String())
			fmt.Println(lib.MachineCode(ctx.DataSection[currentOffset-ctx.DataSectionOffset : len(ctx.DataSection)]))
		}
	}
	fmt.Println(".start:")
	if len(ctx.DataSection) > 0 {
		jmp := asm.JMP(encoding.Uint8(len(ctx.DataSection)))
		fmt.Printf("0x%x: %s\n", 0, jmp.String())
		result_, err := jmp.Encode()
		if err != nil {
			return nil, err
		}
		result = result_
		fmt.Println(lib.MachineCode(result_))
		for _, d := range ctx.DataSection {
			result = append(result, d)
		}
	} else {
		ctx.DataSectionOffset = 0
		ctx.InstructionPointer = 0
	}
	address := uint(ctx.DataSectionOffset + len(ctx.DataSection))
	for _, stmt := range stmts {
		code, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println("\n:: " + stmt.String() + "\n")
		for _, i := range code {
			b, err := i.Encode()
			if err != nil {
				return nil, err
			}
			fmt.Printf("0x%x-0x%x 0x%x: %s\n", address, address+uint(len(b)), ctx.InstructionPointer, i.String())
			address += uint(len(b))
			fmt.Println(lib.MachineCode(b))
			for _, code := range b {
				result = append(result, code)
			}
		}
	}
	fmt.Println()
	return result, nil
}

func CompileIR(stmts []IR) ([]lib.Instruction, error) {
	ctx := NewIRContext()
	for _, stmt := range stmts {
		_, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
	}
	return ctx.GetInstructions(), nil
}

func init() {
	i := []IR{
		MustParseIR(`b = 0; c = 3; d = 3; e = 3; h = 3; z = 300`),
		MustParseIR(`while b != 3 { b = b + 1 }`),
		NewIR_Assignment("a", NewIR_Function(&TFunction{TUint64, []Type{TUint64}, []string{"z"}},
			NewIR_AndThen(
				NewIR_Assignment("b", NewIR_Float64(3.0)),
				NewIR_AndThen(
					NewIR_Assignment("c", NewIR_Cast(NewIR_Variable("z"), TFloat64)),
					NewIR_AndThen(
						NewIR_Assignment("d", NewIR_Mul(NewIR_Variable("b"), NewIR_Variable("c"))),
						NewIR_Return(NewIR_Cast(NewIR_Variable("d"), TUint64))),
				)),
		)),
		NewIR_Assignment("g", NewIR_Call("a", []IRExpression{NewIR_Variable("z")})),
		NewIR_Assignment("g", NewIR_LinuxWrite(NewIR_Uint64(uint64(1)), []uint8("howdy\n"), 6)),
		NewIR_Assignment("g", NewIR_ArrayIndex(NewIR_ByteArray([]uint8("howdy\n")), NewIR_Uint64(2))),
		NewIR_Assignment("g", NewIR_ArrayIndex(NewIR_StaticArray(TUint64,
			[]IRExpression{NewIR_Uint64(0), NewIR_Uint64(1), NewIR_Uint64(4)},
		), NewIR_Uint64(2))),
		MustParseIR(`while g != 4 { g = g + 1 }`),
		NewIR_Assignment("h", NewIR_Struct(
			&TStruct{[]shared.Type{TUint64, TUint64}, []string{"a", "b"}},
			[]IRExpression{NewIR_Uint64(3), NewIR_Uint64(21)},
		)),
		NewIR_Return(NewIR_Variable("g")),
		/*
			NewIR_Assignment("q", NewIR_Float64(2.1415)),
			NewIR_Assignment("q", NewIR_Add(NewIR_Variable("q"), NewIR_Float64(1.5))),
			NewIR_Assignment("i", NewIR_Uint64(0)),
			NewIR_While(NewIR_Not(NewIR_Equals(NewIR_Variable("i"), NewIR_Uint64(5))), NewIR_AndThen(
				NewIR_Assignment("g", NewIR_LinuxWrite(NewIR_Uint64(uint64(1)), []uint8("howdy\n"), 6)),
				NewIR_Assignment("i", NewIR_Add(NewIR_Variable("i"), NewIR_Uint64(1))),
			),
			),
			NewIR_Assignment("j", NewIR_LinuxOpen("/tmp/test.txt", os.O_CREATE|os.O_WRONLY, 0644)),
			NewIR_Assignment("g", NewIR_LinuxWrite(NewIR_Variable("j"), []uint8("howdy, how is it going\n"), 23)),
			NewIR_Return(NewIR_Variable("g")),
		*/
	}
	b, err := Compile(i)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
