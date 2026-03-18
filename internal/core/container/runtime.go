package container

import "context"

type Runtime interface {
	List(ctx context.Context) ([]Instance, error)
	Inspect(ctx context.Context, nameOrId string) (Instance, error)
	Stop(ctx context.Context, nameOrId string) error
	Start(ctx context.Context, nameOrId string) error
}

type Service struct {
	runtime Runtime
}

func NewService(r Runtime) *Service {
	return &Service{
		runtime: r,
	}
}

func (s *Service) ListContainers(ctx context.Context) ([]Instance, error) {
	return s.runtime.List(ctx)
}

func (s *Service) InspectContainer(ctx context.Context, nameOrId string) (Instance, error) {
	return s.runtime.Inspect(ctx, nameOrId)
}

func (s *Service) StopContainer(ctx context.Context, nameOrId string) error {
	return s.runtime.Stop(ctx, nameOrId)
}

func (s *Service) StartContainer(ctx context.Context, nameOrId string) error {
	return s.runtime.Start(ctx, nameOrId)
}
