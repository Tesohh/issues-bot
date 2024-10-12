package slash

import dg "github.com/bwmarrin/discordgo"

func GetOptionMap(i *dg.Interaction) map[string]*dg.ApplicationCommandInteractionDataOption {
	rawOptions := i.ApplicationCommandData().Options

	options := make(map[string]*dg.ApplicationCommandInteractionDataOption, len(rawOptions))
	for _, opt := range rawOptions {
		options[opt.Name] = opt
	}

	return options
}
