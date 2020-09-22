package generator

import (
	"math/rand"
	"strings"
	"time"
)

const base = "AZERTYUIOPMLKJHGFDSQWXCVBN_1234567890-.azertyuiopmlkjhgfdsqwxcvbn"

// MinLength : Minimum Length
const MinLength = 8

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Generate generate random password
func Generate(size int) string {
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteByte(base[rand.Intn(len(base))])
	}

	return sb.String()
}
