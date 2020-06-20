package main

import (
	"fmt"
	"os"
	"strconv"

	"./generator"
)

var size = 16

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if len(os.Args) >= 2 {
		s, err := strconv.Atoi(os.Args[1])
		check(err)
		size = s
	}

	fmt.Printf("** PASSWORD GENERATOR**\npassword length : \033[1;33m%v\033[0m\n", size)
	pwd := generator.Generate(size)
	fmt.Printf("code : \033[1;32m%s\033[0m\n", pwd)
}
