package manifest

import (
	"testing"
)

func TestCompare_NoDrift(t *testing.T) {
	spec := ContainerSpec{
		Name:  "web",
		Image: "nginx:1.25",
		Env:   map[string]string{"PORT": "8080"},
		Ports: []string{"8080:80"},
	}
	result := Compare(spec, spec)
	if result.Drifted {
		t.Errorf("expected no drift, got details: %v", result.Details)
	}
}

func TestCompare_ImageDrift(t *testing.T) {
	spec := ContainerSpec{Name: "web", Image: "nginx:1.25"}
	observed := ContainerSpec{Name: "web", Image: "nginx:1.24"}
	result := Compare(spec, observed)
	if !result.Drifted {
		t.Fatal("expected drift due to image mismatch")
	}
	if len(result.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(result.Details))
	}
}

func TestCompare_EnvDrift_Missing(t *testing.T) {
	spec := ContainerSpec{
		Name:  "app",
		Image: "myapp:latest",
		Env:   map[string]string{"DB_HOST": "localhost"},
	}
	observed := ContainerSpec{
		Name:  "app",
		Image: "myapp:latest",
		Env:   map[string]string{},
	}
	result := Compare(spec, observed)
	if !result.Drifted {
		t.Fatal("expected drift due to missing env var")
	}
}

func TestCompare_EnvDrift_ValueChanged(t *testing.T) {
	spec := ContainerSpec{
		Name:  "app",
		Image: "myapp:latest",
		Env:   map[string]string{"LOG_LEVEL": "info"},
	}
	observed := ContainerSpec{
		Name:  "app",
		Image: "myapp:latest",
		Env:   map[string]string{"LOG_LEVEL": "debug"},
	}
	result := Compare(spec, observed)
	if !result.Drifted {
		t.Fatal("expected drift due to changed env value")
	}
}

func TestCompare_PortDrift(t *testing.T) {
	spec := ContainerSpec{
		Name:  "web",
		Image: "nginx:1.25",
		Ports: []string{"443:443", "80:80"},
	}
	observed := ContainerSpec{
		Name:  "web",
		Image: "nginx:1.25",
		Ports: []string{"80:80"},
	}
	result := Compare(spec, observed)
	if !result.Drifted {
		t.Fatal("expected drift due to missing port")
	}
}
