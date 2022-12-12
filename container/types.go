package container

type ContainerNetwork struct {
	// Host IP address that the container's port is mapped to
	IP string `json:"ip,omitempty"`

	// Port on the container
	// Required: true
	PrivatePort uint16 `json:"privatePort"`

	// Port exposed on the host
	PublicPort uint16 `json:"publicPort,omitempty"`

	// protocoll - tcp / udp
	// Required: false
	Protocol string `json:"protocol"`
}

type ContainerData struct {
	State            string             `json:"state"`
	Status           string             `json:"status"`
	Network          []ContainerNetwork `json:"network"`
	Labels           map[string]string  `json:"labels"`
	ServiceDiscovery bool               `json:"service-discover"`
}

/*
type Container struct {
	Conatiner        map[string]containerData
}
*/

type ByIP []ContainerNetwork

// type SortByPublicPort []containerNetwork
// type SortByPublicProtocol []containerNetwork

func (a ByIP) Len() int           { return len(a) }
func (a ByIP) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIP) Less(i, j int) bool { return a[i].IP < a[j].IP }
