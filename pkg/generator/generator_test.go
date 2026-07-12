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
