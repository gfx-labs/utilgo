package whereami

import (
	"os"
	"strings"
)

var (
	isDockerEnv     bool
	isKubernetesEnv bool
)

func init() {
	isDockerEnv = detectDocker()
	isKubernetesEnv = detectKubernetes()
}

// detectDocker checks if we're running inside a Docker container
func detectDocker() bool {
	// Check for .dockerenv file
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check cgroup for docker
	if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		return strings.Contains(string(data), "docker")
	}

	return false
}

// detectKubernetes checks if we're running inside a Kubernetes cluster
func detectKubernetes() bool {
	// Check for Kubernetes service account
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
		return true
	}

	// Check for Kubernetes environment variables
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	return false
}

// IsDocker returns true if running in a Docker container
func IsDocker() bool {
	return isDockerEnv
}

// IsKubernetes returns true if running in a Kubernetes cluster
func IsKubernetes() bool {
	return isKubernetesEnv
}
