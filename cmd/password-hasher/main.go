package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	password := os.Args[1]
	var cost int
	flag.IntVar(&cost, "cost", 4, "Bcrypt cost")


	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Printf("Failed to generate hash: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(hash))

}
