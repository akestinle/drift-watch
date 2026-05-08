package docker

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	manifest "drift-watch/internal/manifest"
)

func makeInspect(name, image string, env []string, ports nat.PortMap) types.ContainerJSON {
	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Name: "/" + name,
		},
		Config: &container.Config{
			Image: image,
			Env:   env,
		},
		NetworkSettings: &types.NetworkSettings{
			NetworkSettingsBase: types.NetworkSettingsBase{
				Ports: ports,
			},
			Networks: map[string]*network.EndpointSettings{},
		},
	}
}

func TestFromInspect_BasicFields(t *testing.T) {
	info := makeInspect("web", "nginx:latest", []string{"PORT=80"}, nat.PortMap{
		"80/tcp": {},
	})
	cs := fromInspect(info)

	if cs.Name != "web" {
		t.Errorf("expected name %q, got %q", "web", cs.Name)
	}
	if cs.Image != "nginx:latest" {
		t.Errorf("expected image %q, got %q", "nginx:latest", cs.Image)
	}
	if len(cs.Env) != 1 || cs.Env[0] != "PORT=80" {
		t.Errorf("unexpected env: %v", cs.Env)
	}
	if len(cs.Ports) != 1 || !strings.Contains(cs.Ports[0], "80") {
		t.Errorf("unexpected ports: %v", cs.Ports)
	}
}

func TestFromInspect_StripLeadingSlash(t *testing.T) {
	info := makeInspect("myapp", "myapp:v1", nil, nil)
	cs := fromInspect(info)
	if strings.HasPrefix(cs.Name, "/") {
		t.Errorf("name should not start with '/': %q", cs.Name)
	}
}

func TestToManifestContainer(t *testing.T) {
	cs := &ContainerState{
		Name:  "api",
		Image: "api:latest",
		Env:   []string{"DEBUG=true"},
		Ports: []string{"8080/tcp"},
	}
	mc := cs.ToManifestContainer()
	expected := manifest.Container{
		Name:  "api",
		Image: "api:latest",
		Env:   []string{"DEBUG=true"},
		Ports: []string{"8080/tcp"},
	}
	if mc.Name != expected.Name || mc.Image != expected.Image {
		t.Errorf("ToManifestContainer mismatch: got %+v, want %+v", mc, expected)
	}
}
