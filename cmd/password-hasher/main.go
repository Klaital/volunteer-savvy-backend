package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	fmt.Printf("Enter password: ")
	password := ""
	fmt.Scanf("%s", &password)
	cost := 4


	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Printf("Failed to generate hash: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(hash))

}
