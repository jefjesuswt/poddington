package fleet

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type NodeStorer interface {
	Migrate(ctx context.Context) error
	Save(ctx context.Context, n Node) error
	GetAll(ctx context.Context) ([]Node, error)
	GetByID(ctx context.Context, id string) (Node, error)
	Delete(ctx context.Context, id string) error
	UpdateLastSeen(ctx context.Context, id string) error
}

type Service struct {
	storer NodeStorer
}

func NewService(storer NodeStorer) *Service {
	return &Service{
		storer: storer,
	}
}

func (s *Service) Init(ctx context.Context) error {
	return s.storer.Migrate(ctx)
}

func (s *Service) RegisterNode(ctx context.Context, name, address string) (Node, error) {
	nodeID := newSecureString(8)
	token := newSecureString(32)

	node := Node{
		ID:        nodeID,
		Name:      name,
		Address:   address,
		Token:     token,
		CreatedAt: time.Now().UTC(),
		LastSeen:  time.Now().UTC(),
	}

	if err := s.storer.Save(ctx, node); err != nil {
		return Node{}, fmt.Errorf("failed to save node: %w", err)
	}

	return node, nil
}

func (s *Service) ListNodes(ctx context.Context) ([]Node, error) {
	return s.storer.GetAll(ctx)
}

func (s *Service) RemoveNode(ctx context.Context, id string) error {
	return s.storer.Delete(ctx, id)
}

func (s *Service) PingNode(ctx context.Context, id, token string) error {
	node, err := s.storer.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if node.Token != token {
		return fmt.Errorf("unauthorized: invalid token for node: %s", id)
	}

	return s.storer.UpdateLastSeen(ctx, id)
}

func newSecureString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(bytes)
}
