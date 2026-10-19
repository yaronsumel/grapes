package main // import "github.com/yaronsumel/grapes"

import (
	"fmt"
	"os"
)

const (
	version = "0.3.0"
	welcome = `
//      ____ __________ _____  ___  _____
//     / __  / ___/ __  / __ \/ _ \/ ___/
//    / /_/ / /  / /_/ / /_/ /  __(__  )
//    \__, /_/   \__,_/ .___/\___/____/
//   /____/          /_/ v %s // Yaron Sumel [yaronsu@gmail.com]
//
`
)

func main() {
	os.Exit(runCLI("~/.ssh/known_hosts"))
}

func runCLI(knownHostsPath string) (exitCode int) {
	//recover and exit
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "\r\nFatal: %s \n\n", err)
			exitCode = 1
		}
	}()

	fmt.Printf(welcome, version)

	app := newGrape(getInputData())

	app.verifyAction()

	if err := app.run(knownHostsPath, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
