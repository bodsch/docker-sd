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

func findOverwritePort(debug bool, labels map[string]string) (string, string, string) {

	metrics_path_overwrite := ""
	metrics_source_overwrite := ""
	port := ""

	for k, v := range labels {

		if strings.Contains(k, "service-discover") {

			if strings.Contains(k, "service-discover.") {

				if debug {
					fmt.Printf("[DEBUG] |   service-discover overwrites '%s'\n", k)
				}

				if strings.Contains(k, "service-discover.port.") {
					parts := strings.Split(k, ".")
					port = parts[len(parts)-1]

					metrics_path_overwrite = v
				}

				if k == "service-discover.path" {
					metrics_path_overwrite = v
				}
				if k == "service-discover.source" {
					metrics_source_overwrite = v
				}
			}
			continue
		}
	}

	return port, metrics_path_overwrite, metrics_source_overwrite
}

func targetData(debug bool, network []container.ContainerNetwork, labels map[string]string) ([]ServiceDiscover) { // ([1]string, map[string]string) {

	if debug {
		fmt.Printf("[DEBUG] |     - labels: %#v\n", labels)
	}

	port, metrics_path_overwrite, metrics_source_overwrite := findOverwritePort(debug, labels)

	if debug {
		fmt.Printf("[DEBUG] |     - metrics_path_overwrite   : %#s\n", metrics_path_overwrite)
		fmt.Printf("[DEBUG] |     - metrics_source_overwrite : %#s\n", metrics_source_overwrite)
		fmt.Printf("[DEBUG] |     - port                     : %#s\n", port)
	}

	new_labels := make(map[string]string)

	discoveryResult := []ServiceDiscover{}

	var targetArray [1]string

	for idx := 0; idx < len(network); idx++ {
    v := network[idx]

		if debug {
			fmt.Println("[DEBUG] |   +-------------------------------------------------------------------")
			fmt.Printf("[DEBUG] |   |   network: idx: %v, data: %v\n", idx, v)
		}

		private_Port := network[idx].PrivatePort
		metrics_path := ""

		ui64, _ := strconv.ParseUint(port, 10, 64)
		port_as_uint := uint16(ui64)

		if len(utils.KnownMetricsPorts[private_Port]) > 0 {
			metrics_path = utils.KnownMetricsPorts[private_Port]
		}

		// overwrite port and metrics path
		if len(port) > 0 && port_as_uint == private_Port && len(metrics_path_overwrite) > 0 {
			/* convert string to uint16 */
			ui64, _ = strconv.ParseUint(port, 10, 64)
			private_Port = uint16(ui64)

			metrics_path = metrics_path_overwrite
		}
		// ignore all interfaces without metrics path
		if len(metrics_path) == 0 {
			continue
		}

		public_Port  := network[idx].PublicPort

		if debug {
			fmt.Printf("[DEBUG] |   |     -  public port  %v\n", public_Port)
			fmt.Printf("[DEBUG] |   |     -  private port %v\n", private_Port)
			fmt.Printf("[DEBUG] |   |     -  metrics path %v\n", metrics_path)
		}

		new_labels = containerLabels(debug, labels, metrics_path)

		if debug {
			fmt.Printf("[DEBUG] |   |   -  new labels  %v\n", new_labels)
			fmt.Println("[DEBUG] |   +-------------------------------------------------------------------")
		}

		target := fmt.Sprintf("localhost:%s", strconv.FormatInt(int64(public_Port), 10))
		targetArray[0] = target

		discoveryContainer := ServiceDiscover{Targets: targetArray, Labels: new_labels}
		discoveryResult = append(discoveryResult, discoveryContainer)
	}

	return discoveryResult
}

func containerLabels(debug bool, labels map[string]string, metrics_path string) (map[string]string) {

	/** */
	new_labels := make(map[string]string)

	for k, v := range labels {
		// skip some labels
		if strings.Contains(k, "org") || strings.Contains(k, "repository") || strings.Contains(k, "service-discover") {
			if debug {
				fmt.Printf("[DEBUG] |   |     -> ignore label: %s\n", k)
			}
			continue
		}

		if strings.Contains(k, ".") {
			if debug {
				fmt.Printf("[DEBUG] |   |     -> ignore '.' label: %s\n", k)
			}
			continue
		}

		if debug {
			fmt.Printf("[DEBUG] |   |     -> label: %s = %s\n", k, v)
		}

		new_labels[strings.Replace(k, "-", "_", -1)] = v
	}

	new_labels["__metrics_path__"] = metrics_path

	return new_labels
}

func Discover(dockerHost string, debug bool) ([]ServiceDiscover, error) {

	con, _ := container.ListContainer(dockerHost, debug)

	discoveryResult := []ServiceDiscover{}

	for key, value := range con {
		containerDiscover := value.ServiceDiscovery

		if debug {
			fmt.Println("[DEBUG] +-------------------------------------------------------------------")
			fmt.Printf("[DEBUG] | container: %s\n", key)
			fmt.Printf("[DEBUG] |   - service discovery: %v\n", containerDiscover)
		}

		if containerDiscover == true {

			// discoveryContainer := []ServiceDiscover{}

			target_data := targetData(debug, value.Network, value.Labels)

			for _, value := range target_data {
				discoveryContainer := ServiceDiscover{Targets: value.Targets, Labels: value.Labels}
				discoveryResult = append(discoveryResult, discoveryContainer)
			}
		}

		if debug {
			fmt.Println("[DEBUG] +-------------------------------------------------------------------")
		}
	}

	return discoveryResult, nil
}
