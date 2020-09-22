package generator

import "testing"

func TestGenerate(t *testing.T) {

	pwd1 := Generate(MinLength)

	if l := len(pwd1); l != MinLength {
		t.Errorf("generated password's length is not correct : %d", l)
	}

}
