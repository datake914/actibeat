package beater

type Actispy interface {
	getProcessID() (int, error)
	getProcessName() (string, error)
	getWindowName() (string, error)
	getUserName() (string, error)
}
