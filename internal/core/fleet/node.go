package fleet

import (
	"context"
	"time"
)

type Node struct {
	ID string
	Name string
	Address string
	Token string
	CreatedAt time.Time
	LastSeen time.Time
}

type NodeRepository interface {
	Save(ctx context.Context, node Node) error
	GetById(ctx context.Context, id string) (Node, error)
	GetAll(ctx context.Context) ([]Node, error)
	Delete(ctx context.Context, id string) error
	UpdateLastSeen(ctx context.Context, id string) error
}

type Service struct {
	repo NodeRepository
}

func NewService(repo NodeRepository) *Service {
	return &Service{
		repo: repo,
	}
}
