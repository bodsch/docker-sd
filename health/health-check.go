package health

import (
	"encoding/json"
	"github.com/bodsch/container-service-discovery/container"
)

func Check(dockerHost string) (string, error) {

	_, err := container.Ping(dockerHost)

	state := map[string]interface{}{
		"code":   200,
		"health": "okay",
	}

	if err != nil {
		state["code"] = 500
		state["health"] = "failed"
	}

	jsonData, err := json.Marshal(state)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
