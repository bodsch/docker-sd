package full_list

import (
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/utils"
)

func FullList(dockerHost string, debug bool) (string, error) {

	con, _ := container.ListContainer(dockerHost, debug)

	b, err := utils.JSONMarshal(con)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
