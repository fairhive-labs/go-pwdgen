package generator

import (
	"math/rand"
	"strings"
	"time"
)

//Base character reference
const Base = "AZERTYUIOPMLKJHGFDSQWXCVBN_1234567890-.azertyuiopmlkjhgfdsqwxcvbn"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Generate generate random password
func Generate(size int) string {
	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteByte(Base[rand.Intn(len(Base))])
	}

	return sb.String()
}
