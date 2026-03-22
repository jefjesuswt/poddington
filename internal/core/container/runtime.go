package container

import (
	"context"
	"fmt"
)

type Runtime interface {
	List(ctx context.Context, all bool) ([]Instance, error)
	Inspect(ctx context.Context, target string) (Instance, error)
	Stop(ctx context.Context, target string) error
	Start(ctx context.Context, target string) error

	Restart(ctx context.Context, target string) error
	Remove(ctx context.Context, target string, force bool) error
	GetLogs(ctx context.Context, target string) (string, error)
}

type Service struct {
	runtime Runtime
}

func NewService(r Runtime) *Service {
	return &Service{
		runtime: r,
	}
}

func (s *Service) List(ctx context.Context, all bool) ([]Instance, error) {
	return s.runtime.List(ctx, all)
}

func (s *Service) Inspect(ctx context.Context, target string) (Instance, error) {
	return s.runtime.Inspect(ctx, target)
}

func (s *Service) Stop(ctx context.Context, target string) error {
	instance, err := s.runtime.Inspect(ctx, target)
	if err != nil { return err }
	if instance.State == "exited" || instance.State == "stopped" || instance.State == "created" {
		return ErrAlreadyStopped
	}

	return s.runtime.Stop(ctx, target)
}

func (s *Service) Start(ctx context.Context, target string) error {
	instance, err := s.runtime.Inspect(ctx, target)
	if err != nil { return err }
	if instance.State == "running" {
		return ErrAlreadyRunning
	}

	return s.runtime.Start(ctx, target)
}

func (s *Service) Remove(ctx context.Context, target string, force bool) error {
	if !force {
		instance, err := s.runtime.Inspect(ctx, target)
		if err != nil { return err }
		if instance.State == "running" {
			return fmt.Errorf("container is running, stop it first or use --force")
		}
	}
	return s.runtime.Remove(ctx, target, force)
}

func (s *Service) Restart(ctx context.Context, target string) error {
	return s.runtime.Restart(ctx, target)
}

func (s *Service) GetLogs(ctx context.Context, target string) (string, error) {
	return s.runtime.GetLogs(ctx, target)
}
