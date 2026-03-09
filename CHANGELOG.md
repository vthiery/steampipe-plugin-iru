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
