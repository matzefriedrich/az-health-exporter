# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.4.5] - 2026-01-18

* Bumps `github.com/goccy/go-yaml` from 1.19.0 to 1.19.2 [#15](https://github.com/matzefriedrich/az-health-exporter/pull/15)
* Bumps `github.com/matzefriedrich/parsley` from v1.3.0 to v1.3.1
* Bumps `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.20.0 to v1.20.1


## [v0.4.4] - 2025-12-09

* Bumps `github.com/Azure/azure-sdk-for-go/sdk/azidentity` from 1.13.0 to 1.13.1 [#10](https://github.com/matzefriedrich/az-health-exporter/pull/10)
* Bumps `github.com/goccy/go-yaml` from 1.18.0 to 1.19.0 [#11](dependabot/go_modules/github.com/goccy/go-yaml-1.19.0)
* Bumps `github.com/spf13/cobra` from 1.10.1 to 1.10.2 [#12](https://github.com/matzefriedrich/az-health-exporter/pull/12)
* Upgrades Go version to 1.25.5 [#13](https://github.com/matzefriedrich/az-health-exporter/pull/13)


## [v0.4.2] - 2025-11-12

* Bumps `github.com/AzureAD/microsoft-authentication-library-for-go` from 1.5.0 to 1.6.0 [#8](https://github.com/matzefriedrich/az-health-exporter/pull/8)
* Bumps `github.com/Azure/azure-sdk-for-go/sdk/azcore` from 1.19.1 to 1.20.0 [#9](https://github.com/matzefriedrich/az-health-exporter/pull/9)


## [v0.4.1] - 2025-11-04

* Bumps `github.com/matzefriedrich/parsley` from 1.2.1 to 1.3.0 [#7](https://github.com/matzefriedrich/az-health-exporter/pull/7)


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

