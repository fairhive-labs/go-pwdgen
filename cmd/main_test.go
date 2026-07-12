package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

// Controls a length below MinLength is clamped and warned about
func TestGenerateSmallerLength(t *testing.T) {
	var buf bytes.Buffer
	generate(&buf, 5)
	out := buf.String()

	if !strings.Contains(out, "changed to") {
		t.Errorf("expected a warning about the clamped length, got:\n%s", out)
		t.FailNow()
	}
	if want := fmt.Sprintf("Password length : \033[1;33m%v\033[0m", generator.MinLength); !strings.Contains(out, want) {
		t.Errorf("output should report the clamped length %d, got:\n%s", generator.MinLength, out)
		t.FailNow()
	}
}

// Controls a valid length is honored and a password of that length is printed
func TestGenerateValidLength(t *testing.T) {
	const length = 20
	var buf bytes.Buffer
	generate(&buf, length)
	out := buf.String()

	if strings.Contains(out, "changed to") {
		t.Errorf("did not expect a clamp warning for length %d, got:\n%s", length, out)
		t.FailNow()
	}

	// extract the generated password from the "Code : <esc>...<pwd>...<esc>" line
	line := ""
	for l := range strings.SplitSeq(out, "\n") {
		if strings.Contains(l, "Code :") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatalf("output should contain a Code line, got:\n%s", out)
	}

	pwd := strings.TrimSuffix(strings.TrimPrefix(line, "Code : \033[1;32m"), "\033[0m")
	if len(pwd) != length {
		t.Errorf("generated password length is %d, want %d (password=%q)", len(pwd), length, pwd)
		t.FailNow()
	}
}
