
soll

```
  {
    "targets": [
      "localhost:40381"
    ],
    "labels": {
      "container": "workflow-server",
      "application": "workflow-server",
      "service": "coremedia",
      "environment": "stage",
      "source": "actuator",
      "__metrics_path__": "/actuator/prometheus"
    }
  },
  {
    "targets": [
      "localhost:40399"
    ],
    "labels": {
      "container": "workflow-server",
      "application": "workflow-server",
      "service": "coremedia",
      "environment": "stage",
      "source": "metrics",
      "__metrics_path__": "/metrics"
    }
  }

```

ist

```
  {
    "targets": [
      "localhost:8199"
    ],
    "labels": {
      "GIT_BUILD_REF": "cc3c9d2",
      "__metrics_path__": "/metrics",
      "application": "workflow-server",
      "container": "workflow-server",
      "environment": "stage",
      "service-discover": "true",
      "source": "metrics",
      "watchdog": "true"
    }
  }
```
