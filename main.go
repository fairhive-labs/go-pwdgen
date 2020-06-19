package main

import (
	"fmt"
	"math/rand"
	"strings"
)

//Base character reference
const Base = "23456789_-.AZERTYUPMKLJHGFDSQWXCVBNazertyupmkjhgfdsqwxcvbn"

func randomGenerator(size int) string {
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteByte(Base[rand.Intn(len(Base))])
	}

	return sb.String()
}

func main() {
	pwd := randomGenerator(16)
	fmt.Printf("code : \033[1;31m%s\033[0m\n", pwd)
}
