package beater

type ActispyWin32 struct {
}

func (s *ActispyWin32) getProcessID() string {
	return ""
}

func (s *ActispyWin32) getProcessName() string {
	return "processname"
}

func (s *ActispyWin32) getWindowName() string {
	return "windowname"
}

func (s *ActispyWin32) getUserName() string {
	return "username"
}
