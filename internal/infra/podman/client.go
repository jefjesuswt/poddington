package podman

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/jefjesuswt/poddington/internal/core/container"
)

type Client struct {
	// no es un ctx normal, es el ctx custom de podman, con parametros ya inyectados
	// contiene http y socket ya insertado, la api de podman exige usarlo
	podmanCtx context.Context
}

func NewClient() (*Client, error) {
	uid := os.Getuid()
	socketUrl := fmt.Sprintf("unix://run/user/%d/podman/podman.sock", uid)
	ctx := context.Background()

	conn, err := bindings.NewConnection(ctx, socketUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating connection: %w", err)
	}

	return &Client{
		podmanCtx: conn,
	}, nil
}

func (c *Client) List(_ context.Context) ([]container.Instance, error) {

	conts, err := containers.List(c.podmanCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}

	var result []container.Instance

	for _, cont := range conts {
		name := cont.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		result = append(result, container.Instance{
			ID: cont.ID[:12],
			Name: name,
			State: cont.State,
			Image: cont.Image,
		})
	}

	return result, nil
}

func (c *Client) Inspect(_ context.Context, nameOrId string) (container.Instance, error) {

	instance, err := containers.Inspect(c.podmanCtx, nameOrId, nil)
	if err != nil {
		return container.Instance{}, fmt.Errorf("error inspecting container: %w", err)
	}

	return container.Instance{
		ID: instance.ID[:12],
		Name: instance.Name,
		State: instance.State.Status,
		Image: instance.Image,
	}, nil
}

func (c* Client) Start(_ context.Context, nameOrId string) error {
	err := containers.Start(c.podmanCtx, nameOrId, nil)
	if err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}
	return nil
}

func (c* Client) Stop(_ context.Context, nameOrId string) error {
	err := containers.Stop(c.podmanCtx, nameOrId, nil)
	if err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}
	return nil
}
