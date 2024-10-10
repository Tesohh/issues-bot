package main

import (
	"fmt"
	"issues/slash"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var commands = map[string]*slash.Command{}
var registeredCommands = make([]*discordgo.ApplicationCommand, 0)

func registerCommands(session *discordgo.Session) error {
	log.Println("Adding commands...")
	guildid := os.Getenv("DISCORD_GUILD_ID")

	for _, c := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildid, &c.ApplicationCommand)
		if err != nil {
			return fmt.Errorf("Cannot create %s due to %s", c.Name, err.Error())
		}
		registeredCommands = append(registeredCommands, cmd)
	}

	return nil
}
