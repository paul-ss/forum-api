package main

import (
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/server"
)



func main() {
	srv := server.New()
	srv.Run()

	config.Lg("main", "main").Info("Server stopped")
}