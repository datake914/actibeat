package beater

type Actispy interface {
	getProcessID() string
	getProcessName() string
	getWindowName() string
	getUserName() string
}
