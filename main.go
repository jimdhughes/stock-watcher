package main

import (
	"gihtub.com/jimdhughes/stock-watcher/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd := cmd.GetRootCommand()
	cmd.Execute()
}
