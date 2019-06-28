package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"plug/repl"
	"plug/scanner"
)

func main() {
	if len(os.Args) == 1 {
		person, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Plug programming language!\n", person.Username)
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filename := os.Args[1]
		file, err := os.Open(filename)
		defer file.Close()

		if err != nil {
			log.Fatal("unable to read file")
		}

		scanner.Start(file, os.Stdout)
	}
}
