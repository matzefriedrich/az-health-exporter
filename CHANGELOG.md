# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.4.0] - 2025-11-02

### Changed

- Handles HTTP 422 responses from Azure Resource Health API. The monitor keeps tack of unsupported resource types, and skips them during the next health check.


## [v0.3.0] - 2025-11-02

### Added

- Adds the `resource_group` label to the `azure_resource_health_status` metric.`


## [v0.2.0] - 2025-10-31

### Added

- Adds a `Dockerfile` to the project and publishes tagged versions to Docker Hub.
- Adds version metadata to the binary.

## [v0.1.0] - 2025-10-31

### Added

- Initial implementation of the Azure Resource Health exporter.
- Provides a CLI application `az-health-exporter` with `monitor` command.
- Hosts an HTTP server with endpoints to expose Prometheus metrics, resource health status and liveness.
- Adds configuration variables for the configuration of Azure Health API service principal credentials and subscription.

