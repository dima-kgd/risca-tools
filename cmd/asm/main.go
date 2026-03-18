package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/dima-kgd/risca-tools/internal/isa"
)

func main() {
	fmt.Println("RiscA Assembler v.1.0.0")

	ifileName := flag.String("i", "input.asm", "Input file")
	ofileName := flag.String("o", "input.bin", "Output file")
	flag.Usage = func() {
		fmt.Printf("Usage: TODO")
	}
	flag.Parse()

	ifile, err := os.Open(*ifileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer ifile.Close()

	ofile, err := os.OpenFile(*ofileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer ofile.Close()

	scanner := bufio.NewScanner(ifile)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		// tokens, err := isa.Tokenize(line)
		// if err != nil {
		// 	fmt.Printf("Error tokenizing: %v\n", err)
		// 	return
		// }
		// for _, token := range tokens {
		// 	fmt.Print(token, " ")
		// }
		// fmt.Println()
		instruction, err := isa.ParseLine(line)
		if err != nil {
			fmt.Printf("Error parsing: %v\n", err)
			return
		}
		fmt.Printf("0x%04X %s\n", instruction.Pack(), instruction)
		err = binary.Write(ofile, binary.LittleEndian, instruction.Pack())
		if err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		return
	}

}
