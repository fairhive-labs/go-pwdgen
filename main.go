package main

import (
	"fmt"
	"os"
	"strconv"

	"./generator"
)

var size = 16

const minSize = 8

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if len(os.Args) >= 2 {
		s, err := strconv.Atoi(os.Args[1])

		check(err)

		switch {
		case s <= minSize:
			fmt.Printf("\033[1;31mprovided size is less than %v, changed to %v !!!\033[0m\n", minSize, size)
		default:
			size = s
		}
	}

	fmt.Printf("** PASSWORD GENERATOR**\npassword length : \033[1;33m%v\033[0m\n", size)
	pwd := generator.Generate(size)
	fmt.Printf("code : \033[1;32m%s\033[0m\n", pwd)
}
