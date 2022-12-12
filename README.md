# docker-sd

A prometheus service-discovery for docker.

`docker-sd` provides service discovery, which obtains information from a configured dockerd.
The connection can be made either via the local socket or via a remote TCP connection.

## Configuration example

```yaml
rest_api:
  address: "0.0.0.0"
  port: 8088

docker_hosts:
  - host: "tcp://molecule.matrix.lan:2376"
    metrics_ports:
      8080: "/metrics"
```

The information can be retrieved via a REST API.
By default, docker-sd provides an endpoint at `127.0.0.1` on port `8088` for this purpose.

Service Discovery will assign a corresponding target to each service that has an exposed port configured.
The following labels will be attached to each target.

If there are rejecting paths to the export metrics interface, these can be stored accordingly in the configuration.

**Example:**

```json
[
  {
    "targets": [
      "localhost:9090"
    ],
    "labels": {
      "__metrics_path__": "/metrics",
      "application": "prometheus",
      "container": "prometheus",
      "source": "metrics"
    }
  }
]
```

The Service Discover is activated by default for each running container, but can be deactivated via a corresponding label on the container:

```yaml
    labels:
      service-discover.enabled: "false"
```

```yaml
    labels:
      service-discover.enabled: "false"
      service-discover.port.8199: "/metrics"
      service-discover.port.8081: "/actuator/prometheus"
```
