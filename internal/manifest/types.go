package manifest

// ContainerSpec describes a single container's expected configuration.
type ContainerSpec struct {
	Name  string            `yaml:"name"`
	Image string            `yaml:"image"`
	Env   map[string]string `yaml:"env,omitempty"`
	Ports []string          `yaml:"ports,omitempty"`
}

// Manifest represents the desired state of a set of containers.
type Manifest struct {
	Name       string          `yaml:"name"`
	Containers []ContainerSpec `yaml:"containers"`
}

// DriftResult holds the outcome of comparing a running container against its spec.
type DriftResult struct {
	ContainerName string
	Drifted       bool
	Details       []string
}

// Compare checks a running container's observed state against the expected spec
// and returns a DriftResult describing any discrepancies.
func Compare(spec ContainerSpec, observed ContainerSpec) DriftResult {
	result := DriftResult{
		ContainerName: spec.Name,
		Drifted:       false,
	}

	if spec.Image != observed.Image {
		result.Drifted = true
		result.Details = append(result.Details,
			"image mismatch: expected "+spec.Image+", got "+observed.Image)
	}

	for k, expectedVal := range spec.Env {
		if observedVal, ok := observed.Env[k]; !ok {
			result.Drifted = true
			result.Details = append(result.Details, "missing env var: "+k)
		} else if observedVal != expectedVal {
			result.Drifted = true
			result.Details = append(result.Details,
				"env var "+k+" mismatch: expected "+expectedVal+", got "+observedVal)
		}
	}

	specPorts := toSet(spec.Ports)
	observedPorts := toSet(observed.Ports)
	for p := range specPorts {
		if !observedPorts[p] {
			result.Drifted = true
			result.Details = append(result.Details, "missing port binding: "+p)
		}
	}

	return result
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[v] = true
	}
	return s
}
