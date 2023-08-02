package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
	"strconv"
	"strings"
)

// -------------------------

func ListContainer(dockerHost string, debug bool) (map[string]ContainerData, error) {
	ctx := context.Background()

	cli := Client(dockerHost)
	detectedContainer := make(map[string]ContainerData)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("Unable to list containers: %v", err)

		return detectedContainer, err
	}

	if len(containers) > 0 {
		// -> https://pkg.go.dev/github.com/docker/docker@v20.10.21+incompatible/api/types#Container
		for _, container := range containers {

			container_labels := make(map[string]string)
			container_data := ContainerData{}
			ContainerName := container.Names[0]

			if ContainerName[0:1] == "/" {
				ContainerName = ContainerName[1:]
			}

			if debug {
				fmt.Println("[DEBUG] +-------------------------------------------------------------------")
				fmt.Printf("[DEBUG] | container: '%s'\n", ContainerName)
			}
			// enabled at default
			sd_enabled := true

			container_data.Status = container.Status
			container_data.State = container.State

			// define base labels
			container_labels["container"] = ContainerName
			container_labels["application"] = ContainerName

			labels := container.Labels

			for k, v := range labels {
				if debug {
					fmt.Printf("[DEBUG] |    label     : '%s' = '%s'\n", k, v)
				}
				// drop labels
				if k == "maintainer" || k == "owner" || k == "watchdog" || k == "GIT_BUILD_REF" {
					continue
				}

				if strings.Contains(k, "label-schema") || strings.Contains(k, "org.") {
					continue
				}

				if strings.Contains(k, "discover") {
					if k == "service-discover.enabled" || k == "service-discovery" || k == "service-discover" {
						sd_enabled, _ = strconv.ParseBool(v)
					}
					/*
						if sd_enabled == false {
							continue
						} else {
							// TODO
							// custom labels
							// container_data.ServiceDiscovery = sd_enabled
						}
					*/
					// continue
				}
				container_labels[k] = v
			}

			container_data.ServiceDiscovery = sd_enabled
			container_data.Labels = container_labels

			// -> https://pkg.go.dev/github.com/docker/docker@v20.10.21+incompatible/api/types#SummaryNetworkSettings
			ports := container.Ports

			for _, p := range ports {
				// -> https://pkg.go.dev/github.com/docker/docker/api/types#Port
				w := ContainerNetwork{
					IP:          p.IP,
					PrivatePort: p.PrivatePort,
					PublicPort:  p.PublicPort,
					Protocol:    p.Type,
				}
				container_data.Network = append(container_data.Network, w)
			}

			detectedContainer[ContainerName] = container_data
		}

		if debug {
			fmt.Println("[DEBUG] +-------------------------------------------------------------------\n")
		}

	} else {
		log.Println("There are no containers running")
		return detectedContainer, errors.New("There are no containers running")
	}

	/*
		log.Printf("-----------------------------------------------\n")
		// b, err := utils.JSONMarshal(detectedContainer)
		b, err := json.MarshalIndent(detectedContainer, "", "  ")
		if err != nil {
			log.Println(err)
		}

		log.Printf(string(b))
		log.Printf("-----------------------------------------------\n")
	*/

	return detectedContainer, nil
}
