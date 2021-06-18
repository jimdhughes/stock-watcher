package main

import (
	"gihtub.com/jimdhughes/stock-watcher/cmd"
	"github.com/joho/godotenv"
)

func main() {
	initializeEnv()
	cmd := cmd.GetRootCommand()
	cmd.Execute()
}

func initializeEnv() {
	godotenv.Load()
}
