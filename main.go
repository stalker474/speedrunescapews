package main

import (
	"./webapp"
)

func main() {
	//Start the REST server
	webapp.WSApp.Init()
	defer webapp.WSApp.Close()
	webapp.WSApp.StartWS()
}
