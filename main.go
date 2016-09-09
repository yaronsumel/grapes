package main

import (
	"fmt"
	"os"
)

//go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

const (
	version = "0.2.2"
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

	app.run()
}
