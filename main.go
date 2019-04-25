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

	//recover and exit
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\r\nFatal: %s \n\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf(welcome, version)

	app := newGrape(getInputData())

	app.verifyAction()

	app.run("~/.ssh/known_hosts")
}
