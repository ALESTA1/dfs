package config

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func GetRandomKey(m map[string]int) string {
	rand.NewSource(time.Now().UnixNano())
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	// Get a random index
	randomIndex := rand.Intn(len(keys))

	// Return the random key
	return keys[randomIndex]
}

func ResolveHostIp() string {

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
