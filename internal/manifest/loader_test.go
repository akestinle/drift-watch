package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) (dir, name string) {
	t.Helper()
	dir = t.TempDir()
	name = "manifest.yaml"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp manifest: %v", err)
	}
	return dir, name
}

func TestLoad_ValidManifest(t *testing.T) {
	yaml := `
version: "1"
containers:
  - name: web
    image: nginx:1.25
    env:
      PORT: "8080"
    labels:
      app: web
    restartPolicy: always
`
	dir, name := writeTemp(t, yaml)
	l := NewLoader(dir)
	m, err := l.Load(name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(m.Containers))
	}
	c := m.Containers[0]
	if c.Name != "web" {
		t.Errorf("expected name 'web', got %q", c.Name)
	}
	if c.Image != "nginx:1.25" {
		t.Errorf("expected image 'nginx:1.25', got %q", c.Image)
	}
	if c.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", c.Env["PORT"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	l := NewLoader("/nonexistent")
	_, err := l.Load("nope.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MissingName(t *testing.T) {
	yaml := `version: "1"
containers:
  - image: nginx:latest
`
	dir, name := writeTemp(t, yaml)
	_, err := NewLoader(dir).Load(name)
	if err == nil {
		t.Fatal("expected validation error for missing container name")
	}
}

func TestLoad_EmptyContainers(t *testing.T) {
	yaml := `version: "1"
containers: []
`
	dir, name := writeTemp(t, yaml)
	_, err := NewLoader(dir).Load(name)
	if err == nil {
		t.Fatal("expected validation error for empty containers list")
	}
}
