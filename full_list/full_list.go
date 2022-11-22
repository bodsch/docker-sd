package full_list

import (
	"github.com/bodsch/container-service-discovery/container"
	"github.com/bodsch/container-service-discovery/utils"
)

func FullList(dockerHost string) (string, error) {

	con, _ := container.ListContainer(dockerHost)

	b, err := utils.JSONMarshal(con)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
