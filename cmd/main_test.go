package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

// newTheatre returns a theatre wired to in-memory buffers with no delay, so the
// rendering can be asserted on without real sleeps or a terminal.
func newTheatre(animate bool) (theatre, *bytes.Buffer, *bytes.Buffer) {
	var show, out bytes.Buffer
	return theatre{show: &show, out: &out, color: animate, animate: animate}, &show, &out
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
	if strings.Count(out.String(), "\n") != 1 {
		t.Errorf("payload stream must be a single line, got %q", out.String())
	}
	if show.Len() != 0 {
		t.Errorf("bare mode must not emit a decorative show, got %q", show.String())
	}
}

// Controls the piped case (stderr is a TTY, stdout is not): the show renders but
// the payload stays bare so `pwdgen | pbcopy` captures only the password.
func TestRunAnimatedPipedPayloadIsBare(t *testing.T) {
	const length = 24
	tt, show, out := newTheatre(true) // outColor defaults to false = piped stdout
	run(tt, length)

	if pwd := password(out); len(pwd) != length {
		t.Fatalf("payload length is %d, want %d (payload=%q)", len(pwd), length, pwd)
	}
	if strings.Contains(out.String(), "\033") {
		t.Errorf("piped payload must be free of ANSI escapes, got %q", out.String())
	}

	for _, marker := range []string{"PWDGEN", "SECURE FORGE ONLINE", "SEEDING CSPRNG", "PAYLOAD", "100%"} {
		if !strings.Contains(show.String(), marker) {
			t.Errorf("decorative show should contain %q, got:\n%s", marker, show.String())
		}
	}
}

// Controls the fully-interactive case (both streams are TTYs): the payload is
// colored so it reads as part of the interface, yet the raw credential is intact.
func TestRunPayloadColoredOnTTY(t *testing.T) {
	const length = 16
	var show, out bytes.Buffer
	tt := theatre{show: &show, out: &out, color: true, outColor: true, animate: true}
	run(tt, length)

	if !strings.Contains(out.String(), "\033[") {
		t.Errorf("interactive payload should be colored, got %q", out.String())
	}
	// strip the indent and wrapping ANSI codes, confirm the credential is intact
	raw := strings.TrimSpace(out.String())
	raw = strings.TrimPrefix(raw, neon)
	raw = strings.TrimSuffix(raw, reset)
	if len(raw) != length {
		t.Errorf("credential length is %d, want %d (raw=%q)", len(raw), length, raw)
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

// Controls the clamp warning stays free of ANSI escapes when color is off.
func TestRunClampWarningPlain(t *testing.T) {
	tt, show, out := newTheatre(false)
	run(tt, 3)

	if pwd := password(out); len(pwd) != generator.MinLength {
		t.Errorf("payload length is %d, want clamped %d", len(pwd), generator.MinLength)
	}
	if !strings.HasPrefix(show.String(), "!! ") {
		t.Errorf("expected a plain warning, got %q", show.String())
	}
	if strings.Contains(show.String(), "\033") {
		t.Errorf("plain warning must be free of ANSI escapes, got %q", show.String())
	}
}

// Controls the show can play without color: no ANSI escape reaches the stream.
func TestRunAnimatedWithoutColor(t *testing.T) {
	var show, out bytes.Buffer
	run(theatre{show: &show, out: &out, animate: true}, 16)

	if strings.Contains(show.String(), "\033") {
		t.Errorf("uncolored show must be free of ANSI escapes, got %q", show.String())
	}
	if !strings.Contains(show.String(), "FORGING 16-CHAR KEY") {
		t.Errorf("show should announce the key length, got:\n%s", show.String())
	}
}

// Controls the animation really pauses between steps when a delay is set.
func TestPauseHonorsDelay(t *testing.T) {
	const delay = 20 * time.Millisecond
	start := time.Now()
	theatre{delay: delay}.pause()
	if d := time.Since(start); d < delay {
		t.Errorf("paused %v, want at least %v", d, delay)
	}
}

// Controls the progress bar ends full, on its own line.
func TestProgressFillsToFull(t *testing.T) {
	tt, show, _ := newTheatre(false)
	tt.progress()

	if !strings.Contains(show.String(), "["+strings.Repeat("█", 16)+"] 100%") {
		t.Errorf("progress should end on a full bar, got %q", show.String())
	}
	if !strings.HasSuffix(show.String(), "\n") {
		t.Errorf("progress should end with a newline, got %q", show.String())
	}
}

func TestCenter(t *testing.T) {
	tt := []struct {
		name string
		s    string
		w    int
		want string
	}{
		{"even padding", "ab", 6, "  ab  "},
		{"odd padding leans left", "ab", 5, " ab  "},
		{"exact width", "abcd", 4, "abcd"},
		{"wider than field", "abcdef", 4, "abcdef"},
		{"empty", "", 3, "   "},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := center(tc.s, tc.w); got != tc.want {
				t.Errorf("center(%q, %d) = %q, want %q", tc.s, tc.w, got, tc.want)
			}
		})
	}
}

// tempFile returns an open regular file, standing in for a redirected stream.
func tempFile(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

// content returns everything written to f so far.
func content(t *testing.T, f *os.File) string {
	t.Helper()
	b, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestIsTerminal(t *testing.T) {
	if isTerminal(tempFile(t, "regular")) {
		t.Error("a regular file must not be reported as a terminal")
	}

	closed := tempFile(t, "closed")
	closed.Close()
	if isTerminal(closed) {
		t.Error("a file that cannot be stat'ed must not be reported as a terminal")
	}

	// the null device is a character device, like a TTY
	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Skip("no null device:", err)
	}
	defer null.Close()
	if !isTerminal(null) {
		t.Error("a character device should be reported as a terminal")
	}
}

// Controls redirected streams get the bare password and nothing else.
func TestCLIRedirected(t *testing.T) {
	tt := []struct {
		name   string
		args   []string
		length int
	}{
		{"default length", nil, 16},
		{"custom length", []string{"-l", "40"}, 40},
		{"plain", []string{"-plain", "-l", "12"}, 12},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr := tempFile(t, "stdout"), tempFile(t, "stderr")
			if err := cli(tc.args, stdout, stderr); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out := content(t, stdout)
			if pwd := strings.TrimRight(out, "\n"); len(pwd) != tc.length {
				t.Errorf("payload length is %d, want %d (payload=%q)", len(pwd), tc.length, pwd)
			}
			if strings.Contains(out, "\033") {
				t.Errorf("payload must be free of ANSI escapes, got %q", out)
			}
			if show := content(t, stderr); show != "" {
				t.Errorf("redirected stderr must stay empty, got %q", show)
			}
		})
	}
}

// Controls the show plays when stderr is a character device, while a redirected
// stdout still receives the bare password.
func TestCLIAnimatesOnCharDevice(t *testing.T) {
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Skip("no null device:", err)
	}
	defer null.Close()
	if !isTerminal(null) {
		t.Skip("null device is not a character device on this platform")
	}

	stdout := tempFile(t, "stdout")
	start := time.Now()
	if err := cli([]string{"-l", "20"}, stdout, null); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := time.Since(start); d < 18*time.Millisecond {
		t.Errorf("animated run took %v, the show did not play", d)
	}

	out := content(t, stdout)
	if pwd := strings.TrimRight(out, "\n"); len(pwd) != 20 {
		t.Errorf("payload length is %d, want 20 (payload=%q)", len(pwd), pwd)
	}
	if strings.Contains(out, "\033") {
		t.Errorf("redirected payload must be free of ANSI escapes, got %q", out)
	}
}

// Controls -plain silences the show even when stderr is a character device.
func TestCLIPlainSkipsShow(t *testing.T) {
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Skip("no null device:", err)
	}
	defer null.Close()

	stdout := tempFile(t, "stdout")
	start := time.Now()
	if err := cli([]string{"-plain"}, stdout, null); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := time.Since(start); d > 300*time.Millisecond {
		t.Errorf("plain run took %v, the show should not play", d)
	}
	if pwd := strings.TrimRight(content(t, stdout), "\n"); len(pwd) != 16 {
		t.Errorf("payload length is %d, want 16 (payload=%q)", len(pwd), pwd)
	}
}

// Controls bad flags are reported on stderr and no password is produced.
func TestCLIRejectsBadFlags(t *testing.T) {
	for _, args := range [][]string{{"-nope"}, {"-l", "abc"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			stdout, stderr := tempFile(t, "stdout"), tempFile(t, "stderr")
			if err := cli(args, stdout, stderr); err == nil {
				t.Fatal("expected an error")
			}
			if out := content(t, stdout); out != "" {
				t.Errorf("no password should be produced, got %q", out)
			}
			if !strings.Contains(content(t, stderr), "-plain") {
				t.Errorf("usage should be printed on stderr, got %q", content(t, stderr))
			}
		})
	}
}

// Controls -h prints the usage and is not treated as a failure.
func TestCLIHelp(t *testing.T) {
	stdout, stderr := tempFile(t, "stdout"), tempFile(t, "stderr")
	if err := cli([]string{"-h"}, stdout, stderr); err != nil {
		t.Fatalf("help must not be an error, got %v", err)
	}
	if out := content(t, stdout); out != "" {
		t.Errorf("help must not produce a password, got %q", out)
	}
	if !strings.Contains(content(t, stderr), "password length") {
		t.Errorf("usage should be printed on stderr, got %q", content(t, stderr))
	}
}
