package utils

var KnownMetricsPorts = make(map[uint16]string)
var CoremediaServices = []string{}
var CoremediaKnownMetricsPorts = make(map[uint16]string)

func init() {

	KnownMetricsPorts = map[uint16]string{
		/*
			8080: "/metrics", // cadvisor
			8199: "/",        // solr
			9216: "/metrics", // mongodb
			8090: "/metrics", // mgob
			3306: "/metrics", // database
			9090: "/metrics", // prometheus
			9100: "/metrics", // node_exporter
		*/
	}

	CoremediaKnownMetricsPorts = map[uint16]string{
		8081: "/actuator/prometheus",
		8199: "/metrics",
	}

	CoremediaServices = []string{
		"cae-live",
		"cae-preview",
		"studio-client",
		"studio-server",
		"caefeeder-live",
		"caefeeder-preview",
		"content-feeder",
		"workflow-server",
		"master-live-server",
		"content-management-server",
		"replication-live-server",
		"solr-primary",
		"solr-secondary",
		"solr-leader",
		"solr-follower",
		"headless-server",
	}
}
