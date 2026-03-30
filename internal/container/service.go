package container

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) List(ctx context.Context, all bool) ([]Instance, error) {
	return s.repo.List(ctx, all)
}

func (s *Service) Inspect(ctx context.Context, target string) (Instance, error) {
	return s.repo.Inspect(ctx, target)
}

func (s *Service) Stop(ctx context.Context, target string) error {
	instance, err := s.repo.Inspect(ctx, target)
	if err != nil { return err }
	if instance.State == "exited" || instance.State == "stopped" || instance.State == "created" {
		return ErrAlreadyStopped
	}

	return s.repo.Stop(ctx, target)
}

func (s *Service) Start(ctx context.Context, target string) error {
	instance, err := s.repo.Inspect(ctx, target)
	if err != nil { return err }
	if instance.State == "running" {
		return ErrAlreadyRunning
	}

	return s.repo.Start(ctx, target)
}

func (s *Service) Remove(ctx context.Context, target string, force bool) error {
	if !force {
		instance, err := s.repo.Inspect(ctx, target)
		if err != nil { return err }
		if instance.State == "running" {
			return fmt.Errorf("container is running, stop it first or use --force")
		}
	}
	return s.repo.Remove(ctx, target, force)
}

func (s *Service) Restart(ctx context.Context, target string) error {
	return s.repo.Restart(ctx, target)
}

func (s *Service) GetLogs(ctx context.Context, target string) (string, error) {
	return s.repo.GetLogs(ctx, target)
}
