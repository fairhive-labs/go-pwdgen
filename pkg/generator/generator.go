package generator

import (
	"crypto/rand"
	"strings"
)

const base = "^AZERTYUIOPMLKJHGFDSQWXCVBN_#@?1234567890-.!azertyuiopmlkjhgfdsqwxcvbn"

// Password length bounds.
const (
	MinLength = 10
	MaxLength = 1 << 20
)

// threshold is the largest multiple of len(base) that fits in a byte. Random
// bytes greater than or equal to it are rejected so that every character in
// base is drawn with a uniform probability (rejection sampling avoids the
// modulo bias introduced when 256 is not a multiple of len(base)).
const threshold = 256 - (256 % len(base))

// Generate returns a cryptographically secure random password of the requested
// size. If size is outside [MinLength, MaxLength], MinLength is used instead.
func Generate(size int) string {
	// force the minimum length
	if size < MinLength || size > MaxLength {
		size = MinLength
	}

	sb := strings.Builder{}
	sb.Grow(size)

	// buf is filled with random bytes in batches; roughly len(base)/threshold
	// of them are rejected, so refill as needed until size chars are produced.
	buf := make([]byte, size)
	for sb.Len() < size {
		// crypto/rand.Read never returns an error on supported platforms;
		// a failure here means the system CSPRNG is unavailable, so panic.
		if _, err := rand.Read(buf); err != nil {
			panic(err)
		}
		for _, b := range buf {
			if int(b) >= threshold {
				continue // rejected to keep the distribution uniform
			}
			sb.WriteByte(base[int(b)%len(base)])
			if sb.Len() == size {
				break
			}
		}
	}

	return sb.String()
}
