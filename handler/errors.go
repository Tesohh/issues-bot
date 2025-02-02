package handler

import "fmt"

var (
	ErrNotInProject = fmt.Errorf("you're not in a project issue channel")
)
