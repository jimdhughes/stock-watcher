package main

import (
	"gihtub.com/jimdhughes/stock-watcher/cmd"
	"gihtub.com/jimdhughes/stock-watcher/data"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	data.Init("data.bolt")
	cmd := cmd.GetRootCommand()
	cmd.Execute()
}
