package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

// ANSI styling. The interface chrome is green; the payload is bold neon magenta
// so the credential stands out instead of blending into the green.
const (
	reset  = "\033[0m"
	green  = "\033[0;32m"
	bright = "\033[1;32m"
	neon   = "\033[1;95m"
	red    = "\033[1;31m"
)

// theatre carries the two output streams and the animation settings so the
// rendering is fully testable (no real sleeps, no terminal required).
type theatre struct {
	show     io.Writer     // decorative stream (stderr): banner, tasks, progress
	out      io.Writer     // payload stream (stdout): the password
	delay    time.Duration // per-step pause; 0 disables all sleeps
	color    bool          // ANSI green on the decorative stream
	outColor bool          // ANSI green on the payload (only when stdout is a TTY)
	animate  bool          // play the boot sequence + progress bar
}

func main() {
	if err := cli(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		os.Exit(2)
	}
}

// cli parses args and renders a password on the given streams. It is main
// without the process-level side effects, so it can be exercised in tests.
func cli(args []string, stdout, stderr *os.File) error {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)
	length := fs.Int("l", 16, "password length")
	plain := fs.Bool("plain", false, "output the bare password only (no animation), for scripting")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	animate := !*plain && isTerminal(stderr)
	t := theatre{
		show:     stderr,
		out:      stdout,
		color:    animate,
		outColor: animate && isTerminal(stdout),
		animate:  animate,
	}
	if animate {
		t.delay = 18 * time.Millisecond
	}

	run(t, *length)
	return nil
}

// isTerminal reports whether f is attached to a character device (a TTY) rather
// than a pipe or regular file — pure stdlib, no dependency.
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// run generates a password and renders it. The decorative show (when enabled)
// goes to t.show; the password goes to t.out — colored only when stdout is a
// terminal, bare otherwise so pipes stay clean.
func run(t theatre, length int) {
	if length < generator.MinLength {
		t.warn(fmt.Sprintf("requested length %d is below minimum, forced to %d", length, generator.MinLength))
		length = generator.MinLength
	}

	if t.animate {
		t.banner(length)
	}

	pwd := generator.Generate(length)

	if t.animate {
		fmt.Fprintf(t.show, "%s>>> PAYLOAD [%d]:%s\n", t.c(green), length, t.c(reset))
	}
	t.payload(pwd)
	if t.animate {
		fmt.Fprintf(t.show, "%s// NO RIGHT PASSWORD, ONLY BETTER TOOLS%s\n", t.c(green), t.c(reset))
	}
}

// payload writes the credential to the payload stream. When stdout is a TTY it
// is rendered in bright green so it reads as part of the interface; when piped
// or redirected it is written bare so `pwdgen | pbcopy` captures only the
// password.
func (t theatre) payload(pwd string) {
	// On a terminal, indent the credential under the "PAYLOAD" label and paint
	// it neon magenta so it reads as a grouped, highlighted result. When piped
	// or redirected it is written bare (no indent, no color) so scripts capture
	// exactly the password.
	if t.outColor {
		fmt.Fprintf(t.out, "    %s%s%s\n", neon, pwd, reset)
		return
	}
	fmt.Fprintln(t.out, pwd)
}

// banner plays the title box, the task lines and the progress bar on the
// decorative stream.
func (t theatre) banner(length int) {
	const w = 38
	line := strings.Repeat("═", w)
	fmt.Fprintf(t.show, "%s╔%s╗\n", t.c(bright), line)
	fmt.Fprintf(t.show, "║%s║\n", center("-- PWDGEN --", w))
	fmt.Fprintf(t.show, "║%s║\n", center("SECURE FORGE ONLINE", w))
	fmt.Fprintf(t.show, "╚%s╝%s\n", line, t.c(reset))

	for _, label := range []string{
		"SEEDING CSPRNG",
		"REJECTION-SAMPLING ENTROPY",
		fmt.Sprintf("FORGING %d-CHAR KEY", length),
	} {
		t.task(label)
	}

	t.progress()
}

// task prints a single "> LABEL...... [OK]" line, pausing for the animation.
func (t theatre) task(label string) {
	dots := strings.Repeat(".", max(2, 30-len(label)))
	fmt.Fprintf(t.show, "%s> %s%s%s ", t.c(green), label, dots, t.c(reset))
	t.pause()
	fmt.Fprintf(t.show, "%s[OK]%s\n", t.c(bright), t.c(reset))
}

// progress renders a short bar that fills to 100%, redrawing in place with \r.
func (t theatre) progress() {
	const width = 16
	for i := 0; i <= width; i++ {
		bar := strings.Repeat("█", i) + strings.Repeat("░", width-i)
		pct := i * 100 / width
		fmt.Fprintf(t.show, "\r%s[%s] %3d%%%s", t.c(bright), bar, pct, t.c(reset))
		t.pause()
	}
	fmt.Fprintln(t.show)
}

// warn prints a highlighted warning to the decorative stream.
func (t theatre) warn(msg string) {
	if t.color {
		fmt.Fprintf(t.show, "%s!! %s%s\n", red, msg, reset)
		return
	}
	fmt.Fprintf(t.show, "!! %s\n", msg)
}

// c returns s only when color is enabled, so reset codes aren't emitted when
// the decorative stream is plain (keeps test assertions clean).
func (t theatre) c(s string) string {
	if t.color {
		return s
	}
	return ""
}

func (t theatre) pause() {
	if t.delay > 0 {
		time.Sleep(t.delay)
	}
}

// center pads s with spaces so it sits in the middle of a width-w field.
func center(s string, w int) string {
	if len(s) >= w {
		return s
	}
	left := (w - len(s)) / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", w-len(s)-left)
}
