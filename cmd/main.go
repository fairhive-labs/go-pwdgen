package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

func main() {
	l := flag.Int("l", 16, "password length")
	flag.Parse()

	generate(os.Stdout, *l)
}

func generate(w io.Writer, length int) {
	if length < generator.MinLength {
		fmt.Fprintf(w, "\033[1;31mprovided length %v is less than %v, changed to %v !!!\033[0m\n", length, generator.MinLength, generator.MinLength)
		length = generator.MinLength
	}

	fmt.Fprintf(w, "Password length : \033[1;33m%v\033[0m\n", length)
	pwd := generator.Generate(length)
	fmt.Fprintf(w, "Code : \033[1;32m%s\033[0m\n", pwd)
}
