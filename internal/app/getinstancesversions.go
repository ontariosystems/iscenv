package app

import (
	"fmt"
	"strings"

	"github.com/ontariosystems/iscenv/v3/iscenv"
)

// GetInstancesVersions returns list of ISCVersions matching an image for all current ISCInstances
func GetInstancesVersions(image string) (versions iscenv.ISCVersions, err error) {
	instances := GetInstances()
	versions = make(iscenv.ISCVersions, 0, len(instances))
	for _, i := range instances {
		if i.Image == image {
			ai, err := DockerClient.InspectImage(fmt.Sprintf("%s:%s", i.Image, i.Version))
			if err != nil {
				return nil, err
			}

			versions.AddIfMissing(&iscenv.ISCVersion{
				ID:      strings.TrimPrefix(ai.ID, "sha256:"),
				Version: i.Version,
				Created: ai.Created.Unix(),
				Source:  "instances",
			})
		}
	}
	return
}
