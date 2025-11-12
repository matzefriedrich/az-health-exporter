# mobymatze/az-health-exporter

A small service that periodically checks the Azure Health API for a set of configured resources and exposes a simple HTTP service providing health-, status-, and metrics-endpoints.

Images are built and published automatically via GitHub Actions by this repository [github.com/matzefriedrich/az-health-exporter](https://github.com/matzefriedrich/az-health-exporter) when a release tag is pushed. 

The hash of the commit that triggered the CI pipeline is embedded in the binary. For verification, the application prints fully qualified version information along with its banner. Of course, you can build your own image from source.

## How to use this image

Pull the image via `docker pull mobymatze/az-health-exporter:latest`, and run a container; see the following example with a config file mounted at `/config/resources.yaml`:

```
docker run -d \
  --name az-health-exporter \
  -p 8080:8080 \
  -e AZURE_TENANT_ID="<tenant-guid>" \
  -e AZURE_CLIENT_ID="<app-client-id>" \
  -e AZURE_CLIENT_SECRET="<app-client-secret>" \
  -e AZURE_SUBSCRIPTION_ID="<subscription-guid>" \
  -e RESOURCES_CONFIG_FILE="/config/resources.yaml" \
  -e POLL_INTERVAL_SECONDS="60" \
  -v $(pwd)/resources.yaml:/config/resources.yaml:ro \
  mobymatze/az-health-exporter:latest
```

The container listens on port `8080` and serves:

- GET `/health` - A health indicator for the exporter itself.
- GET `/status` - Retrieves the last known health for all configured resources.
- GET `/metrics` - Prometheus metrics

### Docker Compose example

```
services:

  az-health-exporter:
    image: mobymatze/az-health-exporter:latest
    container_name: az-health-exporter

    ports:
      - "8080:8080"

    environment:
      AZURE_CLIENT_ID: ${AZURE_CLIENT_ID}
      AZURE_CLIENT_SECRET: ${AZURE_CLIENT_SECRET}
      AZURE_SUBSCRIPTION_ID: ${AZURE_SUBSCRIPTION_ID}
      AZURE_TENANT_ID: ${AZURE_TENANT_ID}
      POLL_INTERVAL_SECONDS: ${POLL_INTERVAL_SECONDS:-60}
      RESOURCES_CONFIG_FILE: /config/resources.yaml

    volumes:
      - ./resources.yaml:/config/resources.yaml:ro

    restart: unless-stopped
```

## Configuration

The exporter reads configuration from environment variables and a YAML file.

### Required environment variables

- `AZURE_CLIENT_ID` - Azure AD application (client) ID
- `AZURE_CLIENT_SECRET` - Azure AD application client secret
- `AZURE_SUBSCRIPTION_ID` - Azure subscription ID to query
- `AZURE_TENANT_ID` - Azure AD tenant ID
- `RESOURCES_CONFIG_FILE` - Path to the resources YAML file inside the container
- `POLL_INTERVAL_SECONDS` - Optional. Polling interval for health checks (default: 60 seconds)


### Resources YAML format

Provide a YAML file containing the resources list.

Each resource needs its `resource_group`, `name`, and fully-qualified Azure `type` (provider and resource type). For instance:

```
resources:
  - name: "your-domain.de"
    resource_group: "dns-zones"
    type: "Microsoft.Network/dnszones"

  - name: "dumps"
    resource_group: "your-cloud-satellite"
    type: "Microsoft.Storage/storageAccounts"
```

**Unsupported resource types**: The Azure Health API does not provide health information for all resource types. If the exporter encounters a `not supported` response, it traces a warning to `stdout` and skips the resource from polling.


## Metrics

Scrape the `/metrics` endpoint from Prometheus.

Example metrics:

```
# HELP azure_resource_health_status Azure resource health status (1 = healthy, 0 = unhealthy)
# TYPE azure_resource_health_status gauge
azure_resource_health_status{resource_id="...",resource_name="...",resource_type="...",availability_state="Available"} 1

# HELP azure_resource_health_last_check_timestamp Timestamp of last health check
# TYPE azure_resource_health_last_check_timestamp gauge
azure_resource_health_last_check_timestamp{resource_id="...",resource_name="..."} 1.6984752e+09
```

## Source & License

* Source: https://github.com/matzefriedrich/az-health-exporter
* License: [Apache-2.0 license](https://github.com/matzefriedrich/az-health-exporter/blob/main/LICENSE)
