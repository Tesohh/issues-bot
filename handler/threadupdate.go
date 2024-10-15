package handler

import (
	"issues/db"
	"issues/global"
	"issues/slash"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
)

func ThreadUpdate(s *dg.Session, t *dg.ThreadUpdate) {
	var issue db.Issue
	err := global.DB.Find(&issue, "thread_id = ?", t.ID).Error
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if t.ThreadMetadata.Archived &&
		t.ThreadMetadata.ArchiveTimestamp != t.BeforeUpdate.ThreadMetadata.ArchiveTimestamp &&
		issue.IssueStatus > 1 {
		_, err := s.ChannelEdit(t.ID, &dg.ChannelEdit{
			Archived: slash.Ptr(false),
		})
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}
