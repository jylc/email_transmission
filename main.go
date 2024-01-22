package main

import (
	"email_transmission/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("[ERROR] start trans commadn failed, %v\n", err)
	}
}
