package discord

import (
	"fmt"
	"getpack/config"
	"github.com/bwmarrin/discordgo"
	"os"
)

const Logfile = "log.txt"

func SendLog() error {
	cfg := config.Global.Discord

	bot, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return fmt.Errorf("failed to create discord session: %s", err)
	}

	err = bot.Open()
	if err != nil {
		return fmt.Errorf("failed to open bot connection: %s", err)
	}

	file, err := os.Open(Logfile)
	if err != nil {
		return fmt.Errorf("failed to open log file: %s", err)
	}

	_, err = bot.ChannelFileSend(cfg.Channel, Logfile, file)
	if err != nil {
		return fmt.Errorf("failed to send file, %s", err)
	}

	bot.Close()
	return nil
}
