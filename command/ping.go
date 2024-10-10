package command

import (
	"alumix-ceo/slash"

	"github.com/bwmarrin/discordgo"
)

var Ping = slash.Command{
	ApplicationCommand: discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "ping",
	},
	Func: func(s *discordgo.Session, i *discordgo.Interaction) error {
		embed := discordgo.MessageEmbed{
			Title:       "Oi",
			Description: "pinhg",
		}

		return slash.ReplyWithEmbed(s, i, embed)
	},
}
