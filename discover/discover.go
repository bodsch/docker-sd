package discover

import (
	"fmt"
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/utils"
	"reflect"
	"strconv"
	"strings"
)

type ServiceDiscover struct {
	Targets [1]string         `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func UpdateLables(networkInternalPort uint16, labels map[string]string, debug bool) map[string]string {

	new_labels := make(map[string]string)

	if debug {
		fmt.Println("[DEBUG] +-------------------------------------------------------------------")
		fmt.Printf("[DEBUG] | networkInternalPort : '%v'- %v\n", networkInternalPort, reflect.TypeOf(networkInternalPort))
		fmt.Printf("[DEBUG] |  KnownMetricsPorts: '%v'\n", utils.KnownMetricsPorts)
		fmt.Printf("[DEBUG] |  metrics path     : '%v'\n", utils.KnownMetricsPorts[networkInternalPort])
	}

	metrics_path := utils.KnownMetricsPorts[networkInternalPort]

	_actuator := "unknown"
	metrics_path_overwrite := ""
	metrics_source_overwrite := ""

	for k, v := range labels {

		if strings.Contains(k, "org") || strings.Contains(k, "repository") {
			continue
		}

		if strings.Contains(k, "service-discover.") {
			fmt.Printf("[DEBUG] |   service-discover overwrites '%s'\n", k)

			if strings.Contains(k, "service-discover.port.") {
				parts := strings.Split(k, ".")
				port := parts[len(parts)-1]

				if port == strconv.FormatInt(int64(networkInternalPort), 10) {
					metrics_path_overwrite = v
					break
				}
			}

			if k == "service-discover.path" {
				metrics_path_overwrite = v
			}
			if k == "service-discover.source" {
				metrics_source_overwrite = v
			}
		}

		new_labels[k] = v
	}

	if debug {
		fmt.Printf("[DEBUG] | metrics path: '%s' (%v), '%s' (%v)\n", metrics_path, len(metrics_path), metrics_path_overwrite, len(metrics_path_overwrite))
	}

	if len(metrics_path) > 0 || len(metrics_path_overwrite) > 0 {

		if len(metrics_source_overwrite) > 0 {
			new_labels["source"] = metrics_source_overwrite
		} else {

			if strings.Contains(metrics_path, "actuator") {
				_actuator = "actuator"
			}
			if strings.Contains(metrics_path, "metrics") {
				_actuator = "metrics"
			}
			new_labels["source"] = _actuator
		}
		if len(metrics_path_overwrite) > 0 {
			new_labels["__metrics_path__"] = metrics_path_overwrite
		} else {
			new_labels["__metrics_path__"] = metrics_path
		}
	} else {
		// no metrics path defined ...
		new_labels = nil
	}

	if debug {
		fmt.Printf("[DEBUG] | labels: %s\n", new_labels)
	}

	if debug {
		fmt.Println("[DEBUG] +-------------------------------------------------------------------\n")
	}

	return new_labels
}

func Discover(dockerHost string, debug bool) ([]ServiceDiscover, error) {

	con, _ := container.ListContainer(dockerHost, debug)

	discoveryResult := []ServiceDiscover{}

	for key, value := range con {
		if debug {
			fmt.Printf("[DEBUG] container: %s\n", key)
		}
		containerLabels := make(map[string]string)

		containerLabels = value.Labels
		containerNetwork := value.Network
		containerDiscover := value.ServiceDiscovery

		var targetArray [1]string

		if containerDiscover == true {

			var networkPort uint16

			for idx, v := range containerNetwork {
				if debug {
					fmt.Printf("[DEBUG] network: idx: %v, data: %v\n", idx, v)
				}
				networkPort = v.PublicPort
				internalPort := v.PrivatePort

				updated_labels := UpdateLables(internalPort, containerLabels, debug)

				if len(updated_labels) > 0 {

					target := fmt.Sprintf("localhost:%s", strconv.FormatInt(int64(networkPort), 10))
					targetArray[0] = target

					discoveryResult = append(discoveryResult, ServiceDiscover{Targets: targetArray, Labels: updated_labels})
				}
			}
		}

		if debug {
			fmt.Println("[DEBUG] -------------------------------------------------------------------\n")
		}
	}

	return discoveryResult, nil
}
