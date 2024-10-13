package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func cleanFilePaths(filePaths []string) []string {
	var cleanedPaths []string
	for _, filePath := range filePaths {
		cleanedPath := strings.ReplaceAll(filePath, "\\", "/")
		cleanedPath = strings.TrimPrefix(cleanedPath, "data")
		cleanedPaths = append(cleanedPaths, cleanedPath)
	}
	return cleanedPaths
}
func getFilePaths(directory string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	return cleanFilePaths(filePaths), err
}
func Register() {

	filePaths, _ := getFilePaths(Directory)
	storageIP := resolveHostIp()
	clientPort, _ := strconv.Atoi(CLIENT_PORT)
	commandPort, _ := strconv.Atoi(COMMAND_PORT)

	fmt.Println(filePaths)
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
