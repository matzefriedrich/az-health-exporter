# az-health-exporter

A small service that checks Azure Resource Health for a set of Azure resources and exposes:

- A simple HTTP health endpoint (`/health`)
- Prometheus metrics (`/metrics`) about each configured resource’s health

It runs as a CLI app with a `monitor` command that starts an HTTP server and periodically polls Azure for the health of resources.

## Prerequisites

- Go 1.25+ (to build from source)
- An Azure AD application (service principal) with permissions to read resource health for the target subscription; requires the tenant ID and client ID, as well as the client secret.
- The Azure Subscription ID containing the resources you want to monitor.
- A resources configuration file (YAML) listing the resources to check; see the section about the YAML resource format below for details.
- Optional: A Prometheus server to scrape the `/metrics` endpoint

## Build

**Option 1** - install the CLI directly (requires a working Go toolchain):

```bash
go install github.com/matzefriedrich/az-health-exporter/cmd/az-health-exporter@latest
```

The following two options require you to clone the repo first:

```bash
git clone https://github.com/matzefriedrich/az-health-exporter.git
cd az-health-exporter
```

**Option 2** - build a local binary:

```bash
git clone https://github.com/matzefriedrich/az-health-exporter.git
cd az-health-exporter
go build -o az-health-exporter ./cmd/az-health-exporter
```

**Option 3** - build a Docker image
```sh
docker build --rm -t mobymatze/az-health-monitor -f Dockerfile .
```

## Configuration

The monitor reads configuration from environment variables and a YAML file.

### Required environment variables

- `AZURE_CLIENT_ID` – Azure AD application (client) ID
- `AZURE_CLIENT_SECRET` – Azure AD application client secret
- `AZURE_SUBSCRIPTION_ID` – Azure subscription ID to query
- `AZURE_TENANT_ID` – Azure AD tenant ID
- `RESOURCES_CONFIG_FILE` – Path to the resources YAML file
- `POLL_INTERVAL_SECONDS` – Optional. The polling interval for health checks (default: `60` seconds)

### Resources YAML format

Provide a YAML file containing the `resources` list. Each resource needs its `resource_group`, `name`, and fully-qualified Azure `type` (provider, and resource type). Required values can be picked from the response output of `az resource list`. For instance:

```yaml
---
resources:
  - name: "your-domain.de"
    resource_group: "dns-zones"
    type: "Microsoft.Network/dnszones"

  - name: "dumps"
    resource_group: "your-cloud-satellite"
    type: "Microsoft.Storage/storageAccounts"
```

## Usage

After setting environment variables and preparing your YAML file, run the monitor:

```bash
export AZURE_TENANT_ID="<tenant-guid>"
export AZURE_CLIENT_ID="<app-client-id>"
export AZURE_CLIENT_SECRET="<app-client-secret>"
export AZURE_SUBSCRIPTION_ID="<subscription-guid>"
export RESOURCES_CONFIG_FILE="/path/to/resources.yaml"
# Optional
export POLL_INTERVAL_SECONDS=60

az-health-exporter monitor --p 8080
```

### HTTP endpoints

| Endpoint             | Description                                 |
|----------------------|---------------------------------------------|
| **GET** `/health`    | Simple JSON health of the exporter itself   |
| **GET** `/metrics`   | Prometheus metrics for configured resources |

## Metrics

The exporter publishes Prometheus metrics for the configured resources. Scrape the `/metrics` endpoint from Prometheus. Configure your Prometheus server with a job pointing at the host/port you run the exporter on.

---

Copyright 2025 by Matthias Friedrich