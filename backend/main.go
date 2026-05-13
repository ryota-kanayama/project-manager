package main

import (
	"flag"

	"project-manager/config"
	"project-manager/database"
	"project-manager/server"
)

func main() {
	envFile := flag.String("env", "", "path to .env file (optional)")
	flag.Parse()

	config.Init(*envFile)

	if err := database.Init(); err != nil {
		panic(err)
	}
	defer database.Close()

	if err := server.Init(); err != nil {
		panic(err)
	}
}
