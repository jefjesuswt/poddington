package fleet

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Init(ctx context.Context) error {
	return s.repo.Migrate(ctx)
}

func (s *Service) RegisterNode(ctx context.Context, name, address string) (Node, error) {
	nodeID := newSecureString(8)
	token := newSecureString(32)

	node := Node{
		ID: nodeID,
		Name: name,
		Address: address,
		Token: token,
		CreatedAt: time.Now().UTC(),
		LastSeen: time.Now().UTC(),
	}

	if err := s.repo.Save(ctx, node); err != nil {
		return Node{}, fmt.Errorf("failed to save node: %w", err)
	}

	return node, nil
}

func (s *Service) ListNodes(ctx context.Context) ([]Node, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) RemoveNode(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) PingNode(ctx context.Context, id string) error {
	return s.repo.UpdateLastSeen(ctx, id)
}

func newSecureString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(bytes)
}
