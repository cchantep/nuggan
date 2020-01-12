package main

import (
	"fmt"
	"nuggan"
	"os"
)

func main() {
	f, err := os.Open("server.conf")

	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Fails to open configuration file: %s\n",
			err.Error())
		return
	}

	conf, err := nuggan.LoadConfig(f)

	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Fails to load configuration: %s\n",
			err.Error())
		return
	}

	// ---

	fmt.Printf("Configuration: %v\n", conf)

	//TODO
}
