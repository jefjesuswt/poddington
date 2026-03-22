package podman

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
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
