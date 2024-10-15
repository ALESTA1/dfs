package config

import (
	"naming/directree"
	"sync"
)

var Root *directree.Node
var StorageClientPorts = make(map[string]int)
var StorageCommandPorts = make(map[string]int)
var GlobalMutex sync.Mutex
var REGISTRATION_PORT string
var SERVICE_PORT string
