package slash

import (
	"github.com/bwmarrin/discordgo"
)

const EmbedColor = 0xc1121f
const EmbedFooterText = "results provided by Alumix CEO"

func StandardizeEmbed(embed discordgo.MessageEmbed) discordgo.MessageEmbed {
	embed.Color = EmbedColor
	if embed.Footer == nil {
		embed.Footer = &discordgo.MessageEmbedFooter{}
	}
	embed.Footer.Text = EmbedFooterText

	return embed
}
