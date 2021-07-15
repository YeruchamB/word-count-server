package main

func main() {
	var server Server
	server.Initialize()
	server.Run(":8080")
}
