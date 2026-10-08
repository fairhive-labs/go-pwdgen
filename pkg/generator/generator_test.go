package generator

import (
	"strings"
	"testing"
)

// Controls every generated character belongs to the allowed charset
func TestGenerateCharset(t *testing.T) {

	pwd := Generate(4096) // large sample to exercise the whole byte range

	for i, c := range pwd {
		if !strings.ContainsRune(base, c) {
			t.Errorf("TEST charset - character %q at index %d is not part of base", c, i)
			t.FailNow()
		}
	}
}

// Controls default function
func TestGenerateBasic(t *testing.T) {

	pwd := Generate(MinLength)

	if l := len(pwd); l != MinLength {
		t.Errorf("BASIC TEST - generated password's length is not correct : %d", l)
		t.FailNow()
	}
}

// Controls smaller length is changed
func TestGenerateSmallerLength(t *testing.T) {

	pwd := Generate(4)

	if l := len(pwd); l != MinLength {
		t.Errorf("TEST w/ smaller length - generated password's length should be [%d] and not [%d]", MinLength, l)
		t.FailNow()
	}
}

// Controls bigger length is changed
func TestGenerateBiggerLength(t *testing.T) {

	pwd := Generate(MaxLength + 1)
	if l := len(pwd); l != MinLength {
		t.Errorf("incorrect length, got %d, want %d", l, MinLength)
		t.FailNow()
	}
}

// Controls generated password's length are correct
func TestGenerateMultiple(t *testing.T) {

	sizes := []int{
		32, 16, 64, 128, 10, 32, 32, 1024, 2048, 64,
	}
	pwds := make([]string, len(sizes))

	for i, val := range sizes {
		pwds[i] = Generate(val)
		l := len(pwds[i])
		if l != val {
			t.Errorf("TEST w/ multiple size #%d - generated password's length is not correct : expected = %d / generated = %d", i+1, val, l)
			t.FailNow()
		}
	}
}

// Controls generator conflicts
func TestGenerateDetectConflicts(t *testing.T) {

	pwdmap := map[string]bool{}
	max := 10000

	for range max {
		key := Generate(MinLength)
		_, ok := pwdmap[key]
		if ok { // controls key is not already there
			t.Errorf("TEST conflicts - conflicts detected : password [%s] is already present in the map ", key)
			t.FailNow()
		} else {
			pwdmap[key] = true
		}
	}

	if l := len(pwdmap); l != max { // double check with the entire map
		t.Errorf("TEST conflicts - conflicts detected : should be [%d] but got [%d] different passwords", max, l)
		t.FailNow()
	}

}

// Controls the length bounds: both ends are honored, anything outside falls
// back to MinLength.
func TestGenerateBounds(t *testing.T) {
	tt := []struct {
		name string
		size int
		want int
	}{
		{"negative", -1, MinLength},
		{"zero", 0, MinLength},
		{"just below min", MinLength - 1, MinLength},
		{"min", MinLength, MinLength},
		{"just above min", MinLength + 1, MinLength + 1},
		{"max", MaxLength, MaxLength},
		{"just above max", MaxLength + 1, MinLength},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if l := len(Generate(tc.size)); l != tc.want {
				t.Errorf("Generate(%d) has length %d, want %d", tc.size, l, tc.want)
			}
		})
	}
}

// Controls the charset has no duplicate, otherwise some characters would be
// drawn more often than others.
func TestBaseHasNoDuplicate(t *testing.T) {
	seen := map[rune]bool{}
	for _, c := range base {
		if seen[c] {
			t.Errorf("character %q appears more than once in base", c)
		}
		seen[c] = true
	}
}

// Controls the rejection threshold is the largest multiple of len(base) that
// fits in a byte.
func TestThreshold(t *testing.T) {
	if threshold%len(base) != 0 {
		t.Errorf("threshold %d is not a multiple of len(base) %d", threshold, len(base))
	}
	if threshold > 256 || threshold+len(base) <= 256 {
		t.Errorf("threshold %d is not the largest multiple of %d within a byte", threshold, len(base))
	}
}

// Controls every character of base is reachable and none dominates: with
// MaxLength draws each one is expected about 15000 times.
func TestGenerateDistribution(t *testing.T) {
	counts := map[rune]int{}
	for _, c := range Generate(MaxLength) {
		counts[c]++
	}

	want := float64(MaxLength) / float64(len(base))
	for _, c := range base {
		n := float64(counts[c])
		if n < want*0.9 || n > want*1.1 {
			t.Errorf("character %q drawn %d times, want about %.0f", c, counts[c], want)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	for b.Loop() {
		Generate(64)
	}
}
