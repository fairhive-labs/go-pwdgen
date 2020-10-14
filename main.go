package main

import (
	"flag"
	"fmt"
)

func main() {

	l := flag.Int("l", 16, "password length")
	flag.Parse()

	generate(*l)
}

func generate(length int) {
	l := length
	if l < MinLength {
		fmt.Printf("\033[1;31mprovided length is less than %v, changed to %v !!!\033[0m\n", l, MinLength)
		l = MinLength
	}

	fmt.Printf("Password length : \033[1;33m%v\033[0m\n", l)
	pwd := Generate(l)
	fmt.Printf("Code : \033[1;32m%s\033[0m\n", pwd)
}
