package discover

import (
	"fmt"
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/utils"
	"strconv"
	"strings"
)

type ServiceDiscover struct {
	Targets [1]string         `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func UpdateLables(containerNetwork []container.ContainerNetwork, labels map[string]string) map[string]string {

	for _, value := range containerNetwork {

		networkInternalPort := value.PrivatePort

		if val, ok := utils.KnownMetricsPorts[networkInternalPort]; ok {

			_actuator := "unknown"

			if strings.Contains(val, "actuator") {
				_actuator = "actuator"
			}
			if strings.Contains(val, "metrics") {
				_actuator = "metrics"
			}

			labels["__metrics_path__"] = val
			labels["source"] = _actuator

			for k, v := range labels {
				labels[k] = v
			}
		}
	}
	return labels
}

func Discover(dockerHost string) ([]ServiceDiscover, error) {

	con, _ := container.ListContainer(dockerHost)

	discoveryResult := []ServiceDiscover{}

	for _, value := range con {

		containerLabels := make(map[string]string)

		containerLabels = value.Labels
		containerNetwork := value.Network
		containerDiscover := value.ServiceDiscovery

		var targetArray [1]string

		if containerDiscover == true {

			var networkExternalPort uint16

			for _, value := range containerNetwork {
				networkExternalPort = value.PublicPort
			}

			updated_labels := UpdateLables(containerNetwork, containerLabels)

			target := fmt.Sprintf("localhost:%s", strconv.FormatInt(int64(networkExternalPort), 10))
			targetArray[0] = target

			discoveryResult = append(discoveryResult, ServiceDiscover{Targets: targetArray, Labels: updated_labels})
		}
	}

	return discoveryResult, nil
}
