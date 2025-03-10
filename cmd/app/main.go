package main

import (
	"computer-club/internal/server"
)

func main() {
	srv := server.NewServer()
	srv.Run()
}
