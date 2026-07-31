## v0.0.2 [2026-07-31]

_Bug fixes_

- Fixed the `params` column of the `iru_blueprint` table to be decoded as JSON instead of a string. ([#6](https://github.com/vthiery/steampipe-plugin-iru/pull/6))

_Dependencies_

- Bumped `google.golang.org/grpc` from `1.63.2` to `1.79.3`. ([#3](https://github.com/vthiery/steampipe-plugin-iru/pull/3))
- Bumped `github.com/hashicorp/go-getter` from `1.7.5` to `1.8.6`. ([#5](https://github.com/vthiery/steampipe-plugin-iru/pull/5))
- Bumped `go.opentelemetry.io/otel/sdk` from `1.42.0` to `1.43.0`. ([#4](https://github.com/vthiery/steampipe-plugin-iru/pull/4))
- Bumped `github.com/go-jose/go-jose/v4` from `4.1.3` to `4.1.4`. ([#10](https://github.com/vthiery/steampipe-plugin-iru/pull/10))
- Bumped `github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream` from `1.7.7` to `1.7.8`. ([#11](https://github.com/vthiery/steampipe-plugin-iru/pull/11))
- Bumped `golang.org/x/crypto` from `0.49.0` to `0.52.0`. ([#7](https://github.com/vthiery/steampipe-plugin-iru/pull/7))

## v0.0.1 [2026-03-09]

_What's new?_

- Initial release of the plugin with 17 tables:
  - [iru_audit_log](docs/tables/iru_audit_log.md) — list audit log events from the tenant
  - [iru_blueprint](docs/tables/iru_blueprint.md) — list blueprints configured in the tenant
  - [iru_device](docs/tables/iru_device.md) — list and inspect enrolled devices
  - [iru_device_app](docs/tables/iru_device_app.md) — list applications installed on each managed device
  - [iru_device_detail](docs/tables/iru_device_detail.md) — rich device details including MDM status, FileVault, hardware, and Kandji agent info
  - [iru_device_library_item](docs/tables/iru_device_library_item.md) — per-device library item assignment and compliance status
  - [iru_device_parameter](docs/tables/iru_device_parameter.md) — per-device MDM parameter compliance status
  - [iru_library_item](docs/tables/iru_library_item.md) — list library items configured in the tenant
  - [iru_prism_app](docs/tables/iru_prism_app.md) — cross-fleet application inventory from the Prism telemetry endpoint
  - [iru_prism_certificate](docs/tables/iru_prism_certificate.md) — cross-fleet certificate inventory from the Prism telemetry endpoint
  - [iru_prism_filevault](docs/tables/iru_prism_filevault.md) — cross-fleet FileVault encryption status from the Prism telemetry endpoint
  - [iru_self_service_category](docs/tables/iru_self_service_category.md) — list categories defined in the Self Service catalog
  - [iru_tag](docs/tables/iru_tag.md) — list tags configured in the tenant
  - [iru_threat](docs/tables/iru_threat.md) — list threat events detected on managed devices
  - [iru_user](docs/tables/iru_user.md) — list users in the tenant
  - [iru_vulnerability](docs/tables/iru_vulnerability.md) — list CVE-grouped vulnerabilities detected across the tenant
  - [iru_vulnerability_detail](docs/tables/iru_vulnerability_detail.md) — retrieve detailed CVE information including EPSS score and publish dates
