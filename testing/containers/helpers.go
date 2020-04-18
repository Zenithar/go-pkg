package containers

import (
	"strings"

	dockertest "github.com/ory/dockertest/v3"
)

// GetName returns container's name
func GetName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}
