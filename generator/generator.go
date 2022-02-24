package generator

import (
	"math/rand"
	"strings"
	"time"
)

const base = "^AZERTYUIOPMLKJHGFDSQWXCVBN_#@?1234567890-.!azertyuiopmlkjhgfdsqwxcvbn"

// MinLength : Minimum Password's Length
const MinLength = 10

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate random password. If the specified is smaller than the minimum length, default length will be used instead
func Generate(size int) string {

	// force the minimum length
	if size < MinLength {
		size = MinLength
	}

	sb := strings.Builder{}
	for i := 0; i < size; i++ {
		sb.WriteByte(base[rand.Intn(len(base))])
	}

	return sb.String()
}
