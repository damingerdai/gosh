package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func expandPattern(pattern string) string {
	filenames, _ := filepath.Glob(pattern)
	if len(filenames) == 0 {
		return ""
	}
	return strings.Join(filenames, " ")
}

func expandWildcardInCmd(input string) string {

	args := strings.Split(input, " ")

	for i, arg := range args {
		if strings.Contains(arg, "*") || strings.Contains(arg, "?") {
			args[i] = expandPattern(arg)
		}
	}
	fmt.Println(strings.Join(args, " "))
	return strings.Join(args, " ")
}
