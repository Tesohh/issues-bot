package slash

import (
	"github.com/bwmarrin/discordgo"
)

const EmbedColor = 0xffb703

func StandardizeEmbed(embed discordgo.MessageEmbed) discordgo.MessageEmbed {
	embed.Color = EmbedColor

	return embed
}
