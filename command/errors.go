package command

import "fmt"

var (
	ErrProjectAlreadyExists = fmt.Errorf("you already have a project with that prefix")
	ErrNotInIssue           = fmt.Errorf("you are not in an issue thread")
)
