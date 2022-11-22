package container

import (
	"context"
	"encoding/json"
	"fmt"
)

func Ping(dockerHost string) (string, error) {

	cli := *Client(dockerHost)
	ctx := context.Background()

	ping, _ := cli.Ping(ctx)

	data := map[string]interface{}{
		"api_version":  ping.APIVersion,
		"os_type":      ping.OSType,
		"experimental": ping.Experimental,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return "", err
	}

	return string(jsonData), nil
}
