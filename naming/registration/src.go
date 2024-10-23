package registration

import (
	"fmt"
	"naming/config"
	"naming/directree"
	"net/http"
)

func RegistrationService() {
	http.HandleFunc("/register", RegisterHandler)
	config.Root = directree.NewNode("root")
	fmt.Println("Starting registration server on" + config.REGISTRATION_PORT + " ... for Storage registration")
	err := http.ListenAndServe(":"+config.REGISTRATION_PORT, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
