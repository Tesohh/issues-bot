package command

import "fmt"

var (
	ErrProjectAlreadyExists     = fmt.Errorf("you already have a project with that prefix")
	ErrRoleAlreadyExists        = fmt.Errorf("this role is already registered")
	ErrNotInIssue               = fmt.Errorf("you are not in an issue thread")
	ErrRoleIsNotValid           = fmt.Errorf("the role you gave me is not valid")
	ErrIssueIsOldAndICantEditIt = fmt.Errorf("this issue is too old, and i cannot edit it's original message.\nyour changes still got applied.")
)
