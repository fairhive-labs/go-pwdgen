package generator

import (
	"testing"
)

func TestGenerate(t *testing.T) {

	// controls default function
	pwd := Generate(MinLength)

	if l := len(pwd); l != MinLength {
		t.Errorf("BASIC TEST - generated password's length is not correct : %d", l)
		t.FailNow()
	}

	// controls smaller length is changed
	pwd = Generate(4)

	if l := len(pwd); l != MinLength {
		t.Errorf("TEST w/ smaller length - generated password's length should be [%d] and not [%d]", MinLength, l)
		t.FailNow()
	}

	// controls generated password's length are correct

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

	// controls generator conflicts

	pwdmap := map[string]bool{}
	max := 10000

	for i := 0; i < max; i++ {
		key := Generate(MinLength)
		_, ok := pwdmap[key]
		if ok {
			t.Errorf("TEST conflicts - conflicts detected : password [%s] is already present in the map ", key)
			t.FailNow()
		} else {
			pwdmap[key] = true
		}
	}

	if l := len(pwdmap); l != max {
		t.Errorf("TEST conflicts - conflicts detected : should be [%d] but got [%d] different passwords", max, l)
		t.FailNow()
	}

}
