package main

import (
	"log"
)

func main() {
	//Start the REST server
	log.Println("Initializing...")
	WSApp.Init()
	defer WSApp.Close()
	log.Println("Starting webservices...")
	WSApp.StartWS()
}
