# Basics

- [x] Use gorm and sqlite to manage db
- [x] Read mesage
- [x] Mark as done
- [x] Icons to represent issue state
- [x] keep threads always unarchived with the threadupdate event

# Commands

- [x] Change assignees with a slash command
  - [x] Register embed_message_id in the issues tabloe
- [x] Register new roles
- [x] Delete/edit project
  - [x] Move to `project subcommand` structure
- [ ] `list issues` which shows in current project
- [ ] `list projects` which also shows completeness

# Refactoring

- [ ] CLeanup addissue.go
- [x] Move to Issue.KindRole and Issue.PriorityRole

# Automation

- [ ] Read channel mentions from regular messages too
- [ ] Listen for ThreadDelete events and update db accordingly
- [ ] If the thread name changes, change the issue name
