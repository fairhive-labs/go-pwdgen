package main

import (
	"flag"
	"fmt"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

func main() {

	l := flag.Int("l", 16, "password length")
	flag.Parse()

	generate(*l)
}

func generate(length int) {
	l := length
	if l < generator.MinLength {
		fmt.Printf("\033[1;31mprovided length is less than %v, changed to %v !!!\033[0m\n", l, generator.MinLength)
		l = generator.MinLength
	}

	fmt.Printf("Password length : \033[1;33m%v\033[0m\n", l)
	pwd := generator.Generate(l)
	fmt.Printf("Code : \033[1;32m%s\033[0m\n", pwd)
}
