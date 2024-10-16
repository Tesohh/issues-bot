# Basics

- [x] Use gorm and sqlite to manage db
- [x] Read mesage
- [x] Mark as done
- [x] Icons to represent issue state
- [x] keep threads always unarchived with the threadupdate event

# Flexibility

- [ ] Change assignees with a slash command
  - [ ] Register embed_message_id in the issues tabloe
- [ ] Register new roles
- [ ] Delete/edit project

# Refactoring

- [ ] CLeanup addissue.go

# Reactivity

- [ ] Read channel mentions from regular messages too
- [ ] Listen for ThreadDelete events and update db accordingly
- [ ] If the thread name changes, change the issue name
