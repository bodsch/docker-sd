package container

import (
	"github.com/docker/docker/client"
	"log"
)

func Client(dockerHost string) *client.Client {

	// host string, version string, client *http.Client, httpHeaders map[string]string

	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost), client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("Unable to get new docker client: %v", err)
		panic(err)
	}
	defer cli.Close()

	return cli
}
