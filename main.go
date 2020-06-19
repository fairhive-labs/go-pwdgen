package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//Base character reference
const Base = "AZERTYUIOPMLKJHGFDSQWXCVBN_1234567890-.azertyuiopmlkjhgfdsqwxcvbn"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomGenerator generate random password
func RandomGenerator(size int) string {
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteByte(Base[rand.Intn(len(Base))])
	}

	return sb.String()
}

func main() {
	pwd := RandomGenerator(16)
	fmt.Printf("code : \033[1;32m%s\033[0m\n", pwd)
}
