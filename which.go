package main

import (
	"fmt"
	"os/exec"
)

func lookCmmand(cmd string) {
	value := aliasTable[cmd]

	if value != "" {
		fmt.Printf("%s: aliased to %s\n", cmd, value)
		return
	}

	value, err := exec.LookPath(cmd)

	if err == nil {
		fmt.Printf("%s: %s\n", cmd, value)
		return
	}

	fmt.Printf("%s NOT FOUND\n", cmd)
}
