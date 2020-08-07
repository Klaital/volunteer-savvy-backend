package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {

	err := bcrypt.CompareHashAndPassword([]byte(os.Args[1]), []byte(os.Args[2]))
	if err != nil {
		fmt.Printf("%s != %s: %v\n", os.Args[1], os.Args[2], err)
	} else {
		fmt.Printf("%s == %s\n", os.Args[1], os.Args[2])
	}
}
