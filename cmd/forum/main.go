package main

import (
	"github.com/go-park-mail-ru/2020_2_Eternity/configs/config"
	"github.com/paul-ss/forum-api/internal/server"
)



func main() {
	srv := server.New()
	srv.Run()

	config.Lg("main", "main").Info("Server stopped")
}