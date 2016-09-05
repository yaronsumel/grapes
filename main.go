package main

import "fmt"

const (
	version = "0.2"
	welcome = `
//      ____ __________ _____  ___  _____
//     / __  / ___/ __  / __ \/ _ \/ ___/
//    / /_/ / /  / /_/ / /_/ /  __(__  )
//    \__, /_/   \__,_/ .___/\___/____/
//   /____/          /_/ v %s // Yaron Sumel [yaronsu@gmail.com]
//
`
)

var app grape

func main() {
	fmt.Printf(welcome, version)
	app.run()
}
