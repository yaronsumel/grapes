package main

import (
	"fmt"
	"os"
)

func fatal(data string) {
		fmt.Printf("\r\nFatal: %s \n\n" ,data)
		os.Exit(1)
}