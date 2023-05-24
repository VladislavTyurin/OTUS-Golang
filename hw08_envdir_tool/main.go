package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Place your code here.
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir <path_to_en_dir> <command> [agrs]")
		os.Exit(1)
	}
	args := os.Args[1:]
	environment, err := ReadDir(args[0])
	if err != nil {
		log.Fatal(err)
	}
	code := RunCmd(args[1:], environment)
	os.Exit(code)
}
