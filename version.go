package main

import (
	"fmt"
)

var (
	appName = "fsrpl"
	version = "0.2.0"
)

// GetVersion return version string
func GetVersion() string {
	return fmt.Sprintf("%s v%s\n", appName, version)
}

// VersionCmd is commands to display version
type VersionCmd struct{}

// Run is main function
func (v *VersionCmd) Run(opt *Option) error {
	PrintInfof(opt.Stdout, GetVersion())
	return nil
}
