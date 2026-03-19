package container

import "errors"

var (
	ErrAlreadyRunning = errors.New("container is already running!")
	ErrAlreadyStopped = errors.New("container is already stopped!")
)

type Instance struct {
	ID        string
	Name      string
	Image     string
	State     string
	Created   string
	IPAddress string
	Cmd       string
	Ports     []string
	Mounts    []string
}
