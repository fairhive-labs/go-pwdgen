package main

import (
	"fmt"

	"./generator"
)

func main() {
	pwd := generator.Generate(16)
	fmt.Printf("code : \033[1;34m%s\033[0m\n", pwd)
}
