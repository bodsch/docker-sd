package detect

import (
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/utils"
)

type ContainerStates struct {
	State            string `json:"state"`
	Running          bool   `json:"running"`
	ServiceDiscovery bool   `json:"servce-discover"`
}

func Detect(dockerHost string) (string, error) {

	con, _ := container.ListContainer(dockerHost)

	detectedContainers := make(map[string]ContainerStates)

	for key, value := range con {
		running := false

		if value.State == "running" {
			running = true
		}

		w := ContainerStates{
			State:            value.State,
			Running:          running,
			ServiceDiscovery: value.ServiceDiscovery,
		}

		detectedContainers[key] = w
	}

	b, err := utils.JSONMarshal(detectedContainers)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
