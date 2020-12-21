package registry

import (
	"github.com/containrrr/watchtower/pkg/registry/helpers"
	watchtower_types "github.com/containrrr/watchtower/pkg/types"
	ref "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// GetPullOptions creates a struct with all options needed for pulling images from a registry
func GetPullOptions(imageName string) (types.ImagePullOptions, error) {
	auth, err := EncodedAuth(imageName)
	log.Debugf("Got image name: %s", imageName)
	if err != nil {
		return types.ImagePullOptions{}, err
	}

	if auth == "" {
		return types.ImagePullOptions{}, nil
	}
	log.Tracef("Got auth value: %s", auth)

	return types.ImagePullOptions{
		RegistryAuth:  auth,
		PrivilegeFunc: DefaultAuthHandler,
	}, nil
}

// DefaultAuthHandler will be invoked if an AuthConfig is rejected
// It could be used to return a new value for the "X-Registry-Auth" authentication header,
// but there's no point trying again with the same value as used in AuthConfig
func DefaultAuthHandler() (string, error) {
	log.Debug("Authentication request was rejected. Trying again without authentication")
	return "", nil
}

// WarnOnAPIComsumption will return true if the registry is knonw-expected to respond to HTTP HEAD in checking the container digest; will return false if behavior is unknown.
func WarnOnAPIComsumption(container watchtower_types.Container) bool {

	normalizedName, err := ref.ParseNormalizedNamed(container.ImageName())
	if err != nil {
		return false
	}

	container_host, err := helpers.NormalizeRegistry(normalizedName.String())
	if err != nil {
		return false
	}

	if container_host == "index.docker.io" {
		return true
	}

	return false
}
