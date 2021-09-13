package main

import "github.com/compiledpanda/wren-server-go/server"

func main() {
	// TODO #2 call config.Get() and pass in cfg object to server.Setup()
	srv := server.Setup()
	srv.Serve()
}
