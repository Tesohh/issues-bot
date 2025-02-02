package autolist

import (
	"fmt"
	"issues/db"
	"issues/slash"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Embed(title string, defaultPriorityRoleID string, issues []db.Issue) discordgo.MessageEmbed {
	issueStrings := make([]string, 4)

	for _, issue := range issues {
		str := "- "
		str += fmt.Sprintf("<#%s> ", issue.ThreadID)
		if issue.PriorityRoleID != defaultPriorityRoleID {
			str += fmt.Sprintf("<@&%s>", issue.PriorityRoleID)
		}
		issueStrings[issue.IssueStatus] += str + "\n"
	}

	for i := range issueStrings {
		issueStrings[i] = strings.TrimRight(issueStrings[i], "\n")
	}

	description := ""
	if issueStrings[0] != "" {
		description += fmt.Sprintf("**Todo**\n%s\n", issueStrings[0])
	}
	if issueStrings[1] != "" {
		description += fmt.Sprintf("**Doing**\n%s\n", issueStrings[1])
	}
	if issueStrings[2] != "" {
		description += fmt.Sprintf("**Done**\n%s\n", issueStrings[2])
	}
	if issueStrings[3] != "" {
		description += fmt.Sprintf("**Cancelled**\n%s\n", issueStrings[3])
	}

	if description == "" {
		description = "There are no issues here. Get to work!"
	}

	embed := discordgo.MessageEmbed{
		Title:       title,
		Description: description,
	}

	return slash.StandardizeEmbed(embed)
}
