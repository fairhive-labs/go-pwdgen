package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

// ANSI styling.
const (
	reset  = "\033[0m"
	green  = "\033[0;32m"
	bright = "\033[1;32m"
	red    = "\033[1;31m"
)

// tasks are the fake "boot sequence" steps shown before the payload is revealed.
var tasks = []string{
	"INITIALIZING CSPRNG",
	"HARVESTING ENTROPY",
	"BYPASSING WEAK PASSWORDS",
	"FORGING CREDENTIAL",
}

// theatre carries the two output streams and the animation settings so the
// rendering is fully testable (no real sleeps, no terminal required).
type theatre struct {
	show    io.Writer     // decorative stream (stderr): banner, tasks, progress
	out     io.Writer     // payload stream (stdout): the bare password, always
	delay   time.Duration // per-step pause; 0 disables all sleeps
	color   bool          // wrap decorative output in ANSI green
	animate bool          // play the boot sequence + progress bar
}

func main() {
	length := flag.Int("l", 16, "password length")
	plain := flag.Bool("plain", false, "output the bare password only (no animation), for scripting")
	flag.Parse()

	animate := !*plain && isTerminal(os.Stderr)
	t := theatre{
		show:    os.Stderr,
		out:     os.Stdout,
		color:   animate,
		animate: animate,
	}
	if animate {
		t.delay = 45 * time.Millisecond
	}

	run(t, *length)
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
// goes to t.show; the bare password always goes to t.out so pipes stay clean.
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
		fmt.Fprintf(t.show, "\n%s>>> PAYLOAD [%d]:%s\n", t.c(bright), length, t.c(reset))
		fmt.Fprintf(t.show, "%s// THERE IS NO RIGHT PASSWORD, ONLY BETTER TOOLS%s\n", t.c(green), t.c(reset))
	}

	// The password is written raw (never colored) so that `pwdgen | pbcopy`
	// captures exactly the credential and nothing else.
	fmt.Fprintln(t.out, pwd)
}

// banner plays the header, the SYSTEM OVERRIDE box, the task lines and the
// progress bar on the decorative stream.
func (t theatre) banner(length int) {
	fmt.Fprintf(t.show, "%sC:\\> PWDGEN.EXE%s\n", t.c(green), t.c(reset))

	const w = 38
	line := strings.Repeat("═", w)
	fmt.Fprintf(t.show, "%s╔%s╗\n", t.c(bright), line)
	fmt.Fprintf(t.show, "║%s║\n", center("-- SYSTEM OVERRIDE --", w))
	fmt.Fprintf(t.show, "║%s║\n", center("ACCESS GRANTED", w))
	fmt.Fprintf(t.show, "╚%s╝%s\n", line, t.c(reset))

	for _, name := range tasks {
		label := name
		if name == "FORGING CREDENTIAL" {
			label = fmt.Sprintf("%s [%d]", name, length)
		}
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

// progress renders a bar that fills from 0 to 100%, redrawing in place with \r.
func (t theatre) progress() {
	const width = 24
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
