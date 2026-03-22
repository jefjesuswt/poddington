package podman

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func (c *Client) Start(_ context.Context, target string) error {
	if err := containers.Start(c.podmanCtx, target, nil); err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}
	return nil
}

func (c *Client) Stop(_ context.Context, target string) error {
	if err := containers.Stop(c.podmanCtx, target, nil); err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}
	return nil
}

func (c *Client) Restart(_ context.Context, target string) error {
	if err := containers.Restart(c.podmanCtx, target, nil); err != nil {
		return fmt.Errorf("error restarting container: %w", err)
	}
	return nil
}

func (c *Client) Remove(_ context.Context, target string, force bool) error {
	opts := new(containers.RemoveOptions).WithForce(force)
	reps, err := containers.Remove(c.podmanCtx, target, opts);
	if err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}

	for _, r := range reps {
		if r.Err != nil {
			return fmt.Errorf("error removing container: %w", r.Err)
		}
	}
	return nil
}
