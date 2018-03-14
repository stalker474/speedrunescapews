package main

func main() {
	//Start the REST server
	WSApp.Init()
	defer WSApp.Close()
	WSApp.StartWS()
}
