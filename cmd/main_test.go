package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

// newTheatre returns a theatre wired to in-memory buffers with no delay, so the
// rendering can be asserted on without real sleeps or a terminal.
func newTheatre(animate bool) (theatre, *bytes.Buffer, *bytes.Buffer) {
	var show, out bytes.Buffer
	return theatre{show: &show, out: &out, animate: animate}, &show, &out
}

// password returns the single bare line written to the payload stream.
func password(out *bytes.Buffer) string {
	return strings.TrimRight(out.String(), "\n")
}

// Controls the bare (scripting) mode: only the password reaches stdout, clean.
func TestRunBare(t *testing.T) {
	const length = 32
	tt, show, out := newTheatre(false)
	run(tt, length)

	pwd := password(out)
	if len(pwd) != length {
		t.Fatalf("payload length is %d, want %d (payload=%q)", len(pwd), length, pwd)
	}
	if strings.Contains(out.String(), "\033") {
		t.Errorf("payload stream must be free of ANSI escapes, got %q", out.String())
	}
	if strings.ContainsAny(out.String(), "\n") && strings.Count(out.String(), "\n") != 1 {
		t.Errorf("payload stream must be a single line, got %q", out.String())
	}
	if show.Len() != 0 {
		t.Errorf("bare mode must not emit a decorative show, got %q", show.String())
	}
}

// Controls the animated mode: the show is rendered and the payload still clean.
func TestRunAnimated(t *testing.T) {
	const length = 24
	tt, show, out := newTheatre(true)
	run(tt, length)

	if pwd := password(out); len(pwd) != length {
		t.Fatalf("payload length is %d, want %d (payload=%q)", len(pwd), length, pwd)
	}
	if strings.Contains(out.String(), "\033") {
		t.Errorf("payload stream must be free of ANSI escapes, got %q", out.String())
	}

	for _, marker := range []string{"SYSTEM OVERRIDE", "ACCESS GRANTED", "PAYLOAD", "100%"} {
		if !strings.Contains(show.String(), marker) {
			t.Errorf("decorative show should contain %q, got:\n%s", marker, show.String())
		}
	}
}

// Controls that color is applied to the decorative stream only, never leaking
// ANSI escapes into the payload stream (which must stay pipe-clean).
func TestRunColorStaysOffPayload(t *testing.T) {
	const length = 16
	var show, out bytes.Buffer
	tt := theatre{show: &show, out: &out, animate: true, color: true}
	run(tt, length)

	if !strings.Contains(show.String(), "\033[") {
		t.Errorf("colored show should contain ANSI escapes, got:\n%s", show.String())
	}
	if strings.Contains(out.String(), "\033") {
		t.Errorf("payload stream must never contain ANSI escapes, got %q", out.String())
	}
	if pwd := password(&out); len(pwd) != length {
		t.Errorf("payload length is %d, want %d", len(pwd), length)
	}
}

// Controls a length below the minimum is clamped and warned about.
func TestRunClampsShortLength(t *testing.T) {
	tt, show, out := newTheatre(true)
	run(tt, 3)

	if pwd := password(out); len(pwd) != generator.MinLength {
		t.Errorf("payload length is %d, want clamped %d", len(pwd), generator.MinLength)
	}
	if !strings.Contains(show.String(), "minimum") {
		t.Errorf("expected a clamp warning on the decorative stream, got:\n%s", show.String())
	}
}
