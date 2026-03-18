package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/dima-kgd/risca-tools/internal/isa"
)

func main() {
	fmt.Println("RiscA DisAssembler v.1.0.0")
	ifileName := flag.String("i", "input.bin", "Input file")
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

	for {
		var instruction uint16
		err = binary.Read(ifile, binary.LittleEndian, &instruction)
		if err != nil {
			break
		}
		fmt.Printf("0x%04X ", instruction)
		i := isa.Parse(instruction)
		fmt.Println(i)
	}
}
