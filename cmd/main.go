package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DimaGitHahahab/yadro-computer-club/internal/scanner"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: yadro-computer-club <input_file_name>")
	}

	specs, events, err := scanner.Scan(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}
	fmt.Println(specs)
	for _, e := range events {
		fmt.Println(e)
	}
}
