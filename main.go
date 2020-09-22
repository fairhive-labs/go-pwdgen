package main

import (
	"flag"
	"fmt"

	"github.com/trendev/pwdgen/generator"
)

func main() {

	l := flag.Int("l", 16, "password length")
	flag.Parse()

	generate(*l)
}

func generate(length int) {
	l := length
	min := generator.MinLength
	if l < min {
		fmt.Printf("\033[1;31mprovided length is less than %v, changed to %v !!!\033[0m\n", l, min)
		l = min
	}

	fmt.Printf("Password length : \033[1;33m%v\033[0m\n", l)
	pwd := generator.Generate(l)
	fmt.Printf("Code : \033[1;32m%s\033[0m\n", pwd)
}
