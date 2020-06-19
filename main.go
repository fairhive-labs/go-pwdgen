package main

import (
	"fmt"

	"./generator"
)

func main() {
	pwd := generator.Generate(16)
	fmt.Printf("code : \033[1;32m%s\033[0m\n", pwd)
}
