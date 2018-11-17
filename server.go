package main

import (
	"github.com/ScriptArts/DiscordTaskManagementBot/bot"
	"github.com/ScriptArts/DiscordTaskManagementBot/models"
	"github.com/ScriptArts/DiscordTaskManagementBot/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func initialize() error {
	err := utils.LoadEnv()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := initialize()
	if err != nil {
		log.Fatalln(err)
	}

	// discord setting
	discord, err := bot.GetDiscordClient()
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(bot.MentionHandler)

	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}

	models.GetDatabase()

	// システムが終了させられるまで起動し続ける
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
