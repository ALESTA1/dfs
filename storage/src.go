package main

import (
	"fmt"
	"os"
	"time"
)

var CLIENT_PORT string
var COMMAND_PORT string
var REGISTRATION_PORT string
var Directory string

func startCommandServer() {

	Register()

}
func startClientServer() {

}
func main() {

	if len(os.Args) < 4 {
		fmt.Println("Please provide at least 2 ports.")
		return
	}

	CLIENT_PORT = os.Args[1]
	COMMAND_PORT = os.Args[2]
	REGISTRATION_PORT = os.Args[3]
	Directory = os.Args[4]
	os.MkdirAll(Directory, os.ModePerm)

	go startCommandServer()
	go startClientServer()

	for {
		time.Sleep(1000000000 * time.Second)
	}

}
