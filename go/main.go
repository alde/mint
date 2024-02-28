package main

import (
	"fmt"
	"os"
	"os/user"

	"alde.nu/mint/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language REPL!\n", user.Username)

	fmt.Printf("Try out my language by typing in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
