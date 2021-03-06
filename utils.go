package main

import (
	"fmt"
)

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func askForConfirmation(opt *Option) bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		PrintAlertf(opt.Stderr, "cant scan input: %v \n", err)
		return false
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	}
	PrintAlertf(opt.Stdout, "Please type yes or no and then press enter: \n")
	return askForConfirmation(opt)
}
