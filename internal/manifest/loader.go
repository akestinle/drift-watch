package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ContainerSpec represents the desired state of a container as defined in a manifest.
type ContainerSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env"`
	Labels      map[string]string `yaml:"labels"`
	Command     []string          `yaml:"command"`
	RestartPolicy string          `yaml:"restartPolicy"`
}

// Manifest holds the full parsed manifest file.
type Manifest struct {
	Version    string          `yaml:"version"`
	Containers []ContainerSpec `yaml:"containers"`
}

// Loader is responsible for reading and parsing manifest files.
type Loader struct {
	BasePath string
}

// NewLoader creates a new Loader rooted at basePath.
func NewLoader(basePath string) *Loader {
	return &Loader{BasePath: basePath}
}

// Load reads and parses a YAML manifest file by name.
func (l *Loader) Load(filename string) (*Manifest, error) {
	path := filepath.Join(l.BasePath, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %q: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %q: %w", path, err)
	}

	if err := validate(&m); err != nil {
		return nil, fmt.Errorf("invalid manifest %q: %w", path, err)
	}

	return &m, nil
}

// validate performs basic sanity checks on a parsed Manifest.
func validate(m *Manifest) error {
	if len(m.Containers) == 0 {
		return fmt.Errorf("manifest must define at least one container")
	}
	for i, c := range m.Containers {
		if c.Name == "" {
			return fmt.Errorf("container[%d] missing name", i)
		}
		if c.Image == "" {
			return fmt.Errorf("container %q missing image", c.Name)
		}
	}
	return nil
}
