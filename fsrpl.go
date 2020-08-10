package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	au "github.com/logrusorgru/aurora"
)

// EnvDebug is environment variable name for debug
const EnvDebug = "FSRPL_DEBUG"

// Option is command's option
type Option struct {
	Debug          bool
	Stdout, Stderr io.Writer
}

func highlight(str string) au.Value {
	return au.Cyan(str)
}

// Debugf print if debug mode
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// PrintInfof print messages
func PrintInfof(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

// PrintAlertf print messages with alert
func PrintAlertf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprint(w, au.Red(fmt.Sprintf(format, args...)))
}

// Exit codes
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
)

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Version VersionCmd `cmd help:"Output the version number."`
	Copy    CopyCmd    `cmd help:"Copy data from specific firestore path to another firestore path."`
	Dump    DumpCmd    `cmd help:"Dump data to json files."`
	Restore RestoreCmd `cmd help:"Restore data from json files."`
}

func run() int {
	opt := &Option{Debug: cli.Debug, Stdout: os.Stdout, Stderr: os.Stderr}
	ctx := kong.Parse(&cli)
	if cli.Debug {
		os.Setenv(EnvDebug, "1")
	}
	err := ctx.Run(opt)
	if err != nil {
		PrintAlertf(opt.Stderr, "err exit %v \n\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func main() {
	os.Exit(run())
}
