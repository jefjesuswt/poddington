package fleet

import (
	"errors"
	"time"
)

var (
	ErrNodeNotFound = errors.New("node not found")
	ErrNodeAlreadyExists = errors.New("node already exists")
)

type Node struct {
	ID string
	Name string
	Address string
	Token string
	CreatedAt time.Time
	LastSeen time.Time
}
