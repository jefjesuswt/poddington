package container

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/jefjesuswt/poddington/config"
)

type Repository struct {
	client *config.PodmanClient
}

func NewRepository(c *config.PodmanClient) *Repository {
	return &Repository{
		client: c,
	}
}

func (r *Repository) List(_ context.Context, all bool) ([]Instance, error) {
	opts := new(containers.ListOptions).WithAll(all)

	instances, err := containers.List(r.client.Ctx, opts)
	if err != nil {
		return nil, err
	}

	var result []Instance
	for _, instance := range instances {
		name := instance.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		result = append(result, Instance{
			ID:    instance.ID[:12],
			Name:  name,
			State: instance.State,
			Image: instance.Image,
		})
	}

	return result, nil
}

func (r *Repository) Inspect(_ context.Context, target string) (Instance, error) {
	instance, err := containers.Inspect(r.client.Ctx, target, nil)
	if err != nil {
		return Instance{}, err
	}

	// extract ip
	ip := instance.NetworkSettings.IPAddress
	if ip == "" && len(instance.NetworkSettings.Networks) > 0 {
		for _, netw := range instance.NetworkSettings.Networks {
			ip = netw.IPAddress
			break
		}
	}

	// extract command
	cmd := ""
	if len(instance.Config.Cmd) > 0 {
		cmd = strings.Join(instance.Config.Cmd, "")
	}

	// extract port
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

	// extract mount
	var mounts []string
	for _, m := range instance.Mounts {
		mounts = append(mounts, fmt.Sprintf("%s:%s", m.Source, m.Destination))
	}

	return Instance{
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

func (r *Repository) GetLogs(_ context.Context, target string) (string, error) {
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

	err := containers.Logs(r.client.Ctx, target, opts, stdoutChan, stderrChan)
	wg.Wait()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return sb.String(), nil
}

func (r *Repository) Start(_ context.Context, target string) error {
	if err := containers.Start(r.client.Ctx, target, nil); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (r *Repository) Stop(_ context.Context, target string) error {
	if err := containers.Stop(r.client.Ctx, target, nil); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (r *Repository) Restart(_ context.Context, target string) error {
	if err := containers.Restart(r.client.Ctx, target, nil); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (r *Repository) Remove(_ context.Context, target string, force bool) error {
	opts := new(containers.RemoveOptions).WithForce(force)

	reports, err := containers.Remove(r.client.Ctx, target, opts)
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	for _, r := range reports {
		if r.Err != nil {
			return fmt.Errorf("failed to remove container: %w", r.Err)
		}
	}

	return nil
}
