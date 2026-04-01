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

	parser := isa.NewParser()
	scanner := bufio.NewScanner(ifile)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Print(line, "\n")

		skip, err := parser.ParseLine(line)
		if skip {
			continue
		}
		if err != nil {
			fmt.Printf("Line: %d: %v\n", lineNumber, err)
			return
		}
		lineNumber++
	}

	err = parser.ProcessLabels()
	if err != nil {
		fmt.Printf("Error processing labels! %v\n", err)
		return
	}

	for _, instr := range parser.Instructions {
		fmt.Printf("0x%08X 0x%04X %s\n", instr.Address, instr.Pack(), instr)
		err = binary.Write(ofile, binary.LittleEndian, instr.Pack())
		if err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		return
	}

}
