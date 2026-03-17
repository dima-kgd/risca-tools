package main

import (
	"fmt"

	"github.com/dima-kgd/risca-tools/internal/isa"
)

func main() {
	op := isa.GetOpcode(0)
	fmt.Println("RiscA Assembler v.1.0.0")
	fmt.Println(op)
	i := isa.Instruction{Opcode: isa.GetOpcode(0), Func5: 1, Rd: 2, Rs: 3, Ex: 0b01}
	fmt.Println(i)
	i = isa.Instruction{Opcode: isa.GetOpcode(0), Func5: 10, Rd: 2, Rs: 3}
	fmt.Println(i)
	i = isa.Instruction{Opcode: isa.GetOpcode(0), Func5: 13, Rd: 2, Rs: 3}
	fmt.Println(i)
	i = isa.Instruction{Opcode: isa.GetOpcode(1), Func2: 2, Rd: 7, Imm: 0xFF}
	fmt.Println(i)

	asmString := "ADD R1, R2"
	tokens, err := isa.Tokenize(asmString)
	if err != nil {
		fmt.Printf("Error tokenizing: %v\n", err)
		return
	}
	for _, token := range tokens {
		fmt.Print(token, " ")
	}
	fmt.Println()

	instruction, err := isa.ParseLine(asmString)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	fmt.Println(instruction)

}
