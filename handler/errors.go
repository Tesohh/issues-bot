package handler

import "fmt"

var (
	ErrNotInProject         = fmt.Errorf("you're not in a project issue channel")
	ErrProjectHasNoAutolist = fmt.Errorf("this project has no AutoList.\nplease use `/project resendautolist` to enable autolist in this project")
)
