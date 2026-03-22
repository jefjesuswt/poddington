package podman

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/jefjesuswt/poddington/internal/core/container"
)

func (c *Client) List(_ context.Context, all bool) ([]container.Instance, error) {
	opts := new(containers.ListOptions).WithAll(all)

	conts, err := containers.List(c.podmanCtx, opts)
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
			ID:    cont.ID[:12],
			Name:  name,
			State: cont.State,
			Image: cont.Image,
		})
	}
	return result, nil
}

func (c *Client) Inspect(_ context.Context, target string) (container.Instance, error) {
	instance, err := containers.Inspect(c.podmanCtx, target, nil)
	if err != nil {
		return container.Instance{}, fmt.Errorf("error inspecting container: %w", err)
	}

	// ip extraction
	ip := instance.NetworkSettings.IPAddress
	if ip == "" && len(instance.NetworkSettings.Networks) > 0 {
		for _, net := range instance.NetworkSettings.Networks {
			ip = net.IPAddress
			break
		}
	}

	// command extraction
	cmd := ""
	if len(instance.Config.Cmd) > 0 {
		cmd = strings.Join(instance.Config.Cmd, " ")
	}

	// ports extraction
	var ports []string
	for portProto, bindings := range instance.NetworkSettings.Ports {
		if len(bindings) > 0 {
			for _, b := range bindings {
				ports = append(ports, fmt.Sprintf("%s:%s -> %s", b.HostIP, b.HostPort, portProto))
			}
		} else {
			ports = append(ports, string(portProto))
		}
	}

	// mounts extractions
	var mounts []string
	for _, m := range instance.Mounts {
		mounts = append(mounts, fmt.Sprintf("%s:%s", m.Source, m.Destination))
	}

	return container.Instance{
		ID:        instance.ID[:12],
		Name:      instance.Name,
		State:     instance.State.Status,
		Image:     instance.Image,
		Created:   instance.Created.Format("2006-01-02 15:04:05"),
		IPAddress: ip,
		Cmd:       cmd,
		Ports:     ports,
		Mounts:    mounts,
	}, nil
}

func (c *Client) GetLogs(ctx context.Context, target string) (string, error) {
	opts := new(containers.LogOptions).WithStdout(true).WithStderr(true)

	stdoutChan := make(chan string, 100)
	stderrChan := make(chan string, 100)

	var sb strings.Builder
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for line := range stdoutChan {
			sb.WriteString(line)
		}
	}()
	go func() {
		defer wg.Done()
		for line := range stderrChan {
			sb.WriteString(line)
		}
	}()

	err := containers.Logs(c.podmanCtx, target, opts, stdoutChan, stderrChan)
	wg.Wait()

	if err != nil {
		return "", fmt.Errorf("error getting logs: %w", err)
	}

	return sb.String(), nil
}
