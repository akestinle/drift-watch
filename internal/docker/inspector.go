package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	manifest "drift-watch/internal/manifest"
)

// ContainerState holds the runtime state of a container extracted from Docker.
type ContainerState struct {
	Name   string
	Image  string
	Env    []string
	Ports  []string
}

// Inspector wraps the Docker client and provides container inspection helpers.
type Inspector struct {
	cli *client.Client
}

// NewInspector creates an Inspector using the Docker environment defaults.
func NewInspector() (*Inspector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker: failed to create client: %w", err)
	}
	return &Inspector{cli: cli}, nil
}

// Inspect returns the runtime ContainerState for the named container.
func (i *Inspector) Inspect(ctx context.Context, name string) (*ContainerState, error) {
	info, err := i.cli.ContainerInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("docker: inspect %q: %w", name, err)
	}
	return fromInspect(info), nil
}

// ToManifestContainer converts a ContainerState into a manifest.Container so
// it can be compared against a declared manifest.
func (cs *ContainerState) ToManifestContainer() manifest.Container {
	return manifest.Container{
		Name:  cs.Name,
		Image: cs.Image,
		Env:   cs.Env,
		Ports: cs.Ports,
	}
}

// fromInspect maps a Docker inspect response to our ContainerState.
func fromInspect(info types.ContainerJSON) *ContainerState {
	name := strings.TrimPrefix(info.Name, "/")

	var ports []string
	for port := range info.NetworkSettings.Ports {
		ports = append(ports, string(port))
	}

	return &ContainerState{
		Name:  name,
		Image: info.Config.Image,
		Env:   info.Config.Env,
		Ports: ports,
	}
}
