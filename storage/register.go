package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"storage/config"
	"strconv"
	"strings"
)

type RegisterBody struct {
	StorageIP   string   `json:"storage_ip"`
	ClientPort  int      `json:"client_port"`
	CommandPort int      `json:"command_port"`
	Files       []string `json:"files"`
}
type RegisterResponseBody struct {
	Files []string `json:"files"`
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

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdir(1)
	if err == nil {
		return false, nil
	}
	if err == os.ErrNotExist {
		return true, nil
	}
	return false, err
}
func deleteFiles(directory string, files []string) {
	for _, file := range files {
		filePath := filepath.Join(directory, file)

		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("Error removing file %s: %v\n", filePath, err)
			continue
		}
		fmt.Printf("File %s deleted successfully.\n", filePath)

		dirPath := filepath.Dir(filePath)
		for dirPath != directory {

			isEmpty, err := isDirEmpty(dirPath)
			if err != nil {
				fmt.Printf("Error checking directory %s: %v\n", dirPath, err)
				break
			}

			if isEmpty {
				err := os.Remove(dirPath)
				if err != nil {
					fmt.Printf("Error removing directory %s: %v\n", dirPath, err)
					break
				}
				fmt.Printf("Directory %s deleted successfully.\n", dirPath)
			} else {

				break
			}

			dirPath = filepath.Dir(dirPath)
		}
	}
}
func Register() {

	filePaths, _ := getFilePaths(config.Directory)
	storageIP := resolveHostIp()
	clientPort, _ := strconv.Atoi(config.CLIENT_PORT)
	commandPort, _ := strconv.Atoi(config.COMMAND_PORT)

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

	resp, err := http.Post("http://localhost:"+config.REGISTRATION_PORT+"/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var responseBody RegisterResponseBody
		json.NewDecoder(resp.Body).Decode(&responseBody)
		fmt.Println(responseBody.Files)
		deleteFiles(config.Directory, responseBody.Files)
		fmt.Println("Successfully registered.")
	} else {
		fmt.Println("Failed to register, status code:", resp.StatusCode)
	}
}
