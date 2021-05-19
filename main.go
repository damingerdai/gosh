package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	initialize()

	stdin := bufio.NewReader(os.Stdin)

	for {
		showPrompt()

		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)

		err := executeInput(input)
		if err != nil {
			log.Println(err)
		}
	}
}

func executeInput(input string) error {
	if strings.HasPrefix(input, `\`) {
		input = input[1:]
	} else {
		input = expandAlias(input)
	}

	input = expandWildcardInCmd(input)
	input = os.ExpandEnv(input)

	args := parseArgs(strings.Trim(input, " "))

	if args[0] == "which" {
		for _, cmd := range args[1:] {
			lookCmmand(cmd)
		}
		return nil
	}

	if args[0] == "alias" {
		kv := strings.Split(args[1], "=")

		key, val := kv[0], strings.Trim(kv[1], "'")
		setAlias(key, val)

		return nil
	}

	if args[0] == "unalias" {
		key := args[1]

		unsetAlias(key)

		return nil
	}

	if args[0] == "export" {
		kv := strings.Split(args[1], "=")

		key, val := kv[0], kv[1]

		err := os.Setenv(key, val)
		return err
	}

	if args[0] == "unset" {
		err := os.Unsetenv(args[1])

		return err
	}

	if args[0] == "cd" {
		err := os.Chdir(args[1])

		return err
	}

	if args[0] == "exit" || args[0] == "exit:" {
		os.Exit(0)
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func initialize() {
	homeDir, _ := os.UserHomeDir()

	file, err := os.Open(fmt.Sprintf("%s/.goshrc", homeDir))

	if err != nil {
		return
	}

	goshrcReader := bufio.NewReader(file)

	for {
		input, err := goshrcReader.ReadString('\n')

		if err == io.EOF {
			return
		}

		input = strings.TrimSpace(input)

		executeInput(input)
	}
}
