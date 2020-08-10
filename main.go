package main

import "os"

func main() {
	cli := CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	code := cli.Run(os.Args)
	os.Exit(code)
}
