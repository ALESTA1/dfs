package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var CLIENT_PORT string
var COMMAND_PORT string
var REGISTRATION_PORT string
var Directory string

type RegisterBody struct {
	StorageIP   string   `json:"storage_ip"`
	ClientPort  int      `json:"client_port"`
	CommandPort int      `json:"command_port"`
	Files       []string `json:"files"`
}

func resolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {

			ip := networkIp.IP.String()

			fmt.Println("Resolved Host IP: " + ip)

			return ip
		}
	}
	return ""
}

func register() {

	//process the data directory to get file paths
	var filePaths = []string{"file1.txt", "file2.txt", "file3.txt"}
	storageIP := resolveHostIp()
	clientPort, _ := strconv.Atoi(CLIENT_PORT)
	commandPort, _ := strconv.Atoi(COMMAND_PORT)

	fmt.Println(storageIP)

	body := RegisterBody{
		StorageIP:   storageIP,
		ClientPort:  clientPort,
		CommandPort: commandPort,
		Files:       filePaths,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:"+REGISTRATION_PORT+"/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		//process files here
		fmt.Println("Successfully registered.")
	} else {
		fmt.Println("Failed to register, status code:", resp.StatusCode)
	}
}

func startCommandServer() {

	register()

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
