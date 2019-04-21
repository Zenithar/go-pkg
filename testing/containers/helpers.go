package containers

import (
	"strings"
	
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

// GetName returns container's name
func GetName(container *dockertest.Resource) string {
	return strings.TrimPrefix(container.Container.Name, "/")
}