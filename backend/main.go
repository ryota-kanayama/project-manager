package main

import (
	"flag"

	"project-manager/config"
	"project-manager/database"
	"project-manager/logger"
	"project-manager/server"
)

func main() {
	envFile := flag.String("env", "", "path to .env file (optional)")
	flag.Parse()

	config.Init(*envFile)
	logger.Init()

	if err := database.Init(); err != nil {
		panic(err)
	}
	defer database.Close()

	if err := server.Init(); err != nil {
		panic(err)
	}
}
