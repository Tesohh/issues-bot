package slash

import (
	"github.com/bwmarrin/discordgo"
)

func ReplyWithEmbed(s *discordgo.Session, i *discordgo.Interaction, embed discordgo.MessageEmbed) error {
	embed = StandardizeEmbed(embed)

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
