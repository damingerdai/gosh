package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var currentCmd *exec.Cmd

func main() {
	initialize()

	signalCh := make(chan os.Signal)

	signal.Notify(signalCh, syscall.SIGINT)

	stdin := bufio.NewReader(os.Stdin)

	handleSignals := func() {
		for {
			sig := <-signalCh

			if currentCmd != nil {
				currentCmd.Process.Signal(sig)
			} else {
				fmt.Println()
				showPrompt()
			}

		}
	}
	go handleSignals()

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

	shouldRunInBackground := false

	inputStream := os.Stdin

	args := parseArgs(strings.Trim(input, " "))

	if args[len(args)-1] == "&" {
		shouldRunInBackground = true
		inputStream = nil
		args = args[:len(args)-1]
	}

	if len(args) > 2 && args[len(args)-2] == "<" {
		filename := args[len(args)-1]

		file, err := os.Open(filename)

		if err != nil {
			return err
		}

		inputStream = file

		args = args[:len(args)-2]

	}

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

	cmd.Stdin = inputStream
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if shouldRunInBackground {
		err := cmd.Start()
		return err
	}

	currentCmd = cmd
	err := cmd.Run()

	currentCmd = nil

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
