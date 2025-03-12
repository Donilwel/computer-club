package main

import (
	"computer-club/internal/server"
)

// @title           API Computer-club
// @version         1.0
// @description     This is a sample API with Swagger
// @host           localhost:8080
// @BasePath       /
func main() {
	srv := server.NewServer()
	srv.Run()
}
