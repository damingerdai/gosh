package main

import "strings"

func parseArgs(input string) []string {
	if strings.HasPrefix(input, "alias") {
		return strings.SplitN(input, " ", 2)
	}

	return strings.Split(input, " ")
}
