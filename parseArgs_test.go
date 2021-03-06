package main

import "testing"

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestParseArgs(t *testing.T) {
	input := "ls -l -a"
	ans := []string{"ls", "-l", "-a"}

	if !equal(parseArgs(input), ans) {
		t.Fail()
	}

	input = "alias gst='git status'"
	ans = []string{"alias", "gst='git status'"}

	if !equal(parseArgs(input), ans) {
		t.Fail()
	}
}
