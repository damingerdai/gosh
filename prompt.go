package main

import (
	"fmt"
	"os"
	"os/user"
)

func showPrompt() {
	u, _ := user.Current()
	host, _ := os.Hostname()
	wd, _ := os.Getwd()

	userAndHost := blue(fmt.Sprintf("%s@%s", u.Username, host))
	wd = yellowWithBlueBG(wd)

	fmt.Printf("%s %s > ", userAndHost, wd)
}
