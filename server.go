package main

import (
	"github.com/ScriptArts/DiscordTaskManagementBot/utils"
	"log"
)

func initialize() error {
	err := utils.LoadEnv()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	initialize()
	log.Println("Hello Bot")
}
