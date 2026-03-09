# Iru Plugin for Steampipe

Use SQL to query devices, users, blueprints, threats, library items, and more from [Iru](https://iru.com) (formerly Kandji), the Apple endpoint management platform.

- **[Get started →](https://github.com/vthiery/steampipe-plugin-iru)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-iru/issues)

## Quick start

### Install

Install the plugin with [Steampipe](https://steampipe.io):

```sh
steampipe plugin install ghcr.io/vthiery/iru
```

### Configure

Copy the sample config and set your API credentials:

```sh
cp config/iru.spc ~/.steampipe/config/iru.spc
```

Edit `~/.steampipe/config/iru.spc`:

```hcl
connection "iru" {
  plugin = "local/iru"

  # The full API URL for your Iru tenant.
  # US tenants: https://<subdomain>.api.kandji.io
  # EU tenants: https://<subdomain>.api.eu.kandji.io
  url = "https://acme.api.eu.kandji.io"

  # API token generated in Settings > Access > API Token.
  token = "YOUR_API_TOKEN"
}
```

### Run a query

```shell
steampipe query
```

List all enrolled Mac devices on an outdated OS:

```sql
select
  device_name,
  serial_number,
  os_version,
  last_check_in
from
  iru_device
where
  platform = 'Mac'
  and os_version < '14.0';
```

## Tables

### Device inventory

| Table | Description |
|-------|-------------|
| [iru_device](iru/table_iru_device.go) | All devices enrolled in your Iru tenant. |
| [iru_device_detail](iru/table_iru_device_detail.go) | Rich per-device details: MDM, FileVault, hardware, Kandji agent, activation lock, and ADE status. |
| [iru_device_app](iru/table_iru_device_app.go) | Applications installed on each managed device. |
| [iru_device_library_item](iru/table_iru_device_library_item.go) | Per-device library item assignment and compliance status. |
| [iru_device_parameter](iru/table_iru_device_parameter.go) | Per-device MDM parameter compliance status. |

### Users, blueprints & tags

| Table | Description |
|-------|-------------|
| [iru_user](iru/table_iru_user.go) | Users in your Iru tenant. |
| [iru_blueprint](iru/table_iru_blueprint.go) | Blueprints configured in your Iru tenant. |
| [iru_tag](iru/table_iru_tag.go) | Tags configured in your Iru tenant. |

### Library, threats & audit

| Table | Description |
|-------|-------------|
| [iru_library_item](iru/table_iru_library_item.go) | Custom app library items in your Iru tenant. |
| [iru_threat](iru/table_iru_threat.go) | Threat events detected on managed devices. |
| [iru_audit_log](iru/table_iru_audit_log.go) | Audit log events from your Iru tenant. |

### Prism telemetry

These tables require a **Prism-enabled API token** (enabled separately in Iru Settings → Access).

| Table | Description |
|-------|-------------|
| [iru_prism_app](iru/table_iru_prism_app.go) | Cross-fleet application inventory from the Prism telemetry endpoint. |
| [iru_prism_certificate](iru/table_iru_prism_certificate.go) | Cross-fleet certificate inventory from the Prism telemetry endpoint. |
| [iru_prism_filevault](iru/table_iru_prism_filevault.go) | Cross-fleet FileVault encryption status from the Prism telemetry endpoint. |

### Self Service

| Table | Description |
|-------|-------------|
| [iru_self_service_category](iru/table_iru_self_service_category.go) | Categories defined in the Iru Self Service catalog. |

### Vulnerability Management

These tables require the **Vulnerability Management add-on** to be enabled on your Iru tenant. The endpoints return 404 without the add-on.

| Table | Description |
|-------|-------------|
| [iru_vulnerability](iru/table_iru_vulnerability_summary.go) | CVE-grouped vulnerabilities detected across the tenant. |
| [iru_vulnerability_detail](iru/table_iru_software_vulnerability.go) | Full CVE detail (description, EPSS score, publish dates) for a specific CVE ID. |

> **Note — untested tables:** Several tables could not be fully tested because the API token used during development lacks the required permission scopes. These tables are correctly implemented against the [Iru API documentation](https://api-docs.kandji.io) and will work once the appropriate scopes are enabled on your token:
>
> - `iru_blueprint`, `iru_tag` — requires the *Blueprints* / *Tags* API scope
> - `iru_library_item`, `iru_threat`, `iru_audit_log` — requires the respective *Library Items* / *Threats* / *Audit Log* scopes
> - `iru_device_library_item`, `iru_device_parameter` — requires the *Device Information* scope
> - `iru_prism_app`, `iru_prism_certificate`, `iru_prism_filevault` — requires a **Prism-enabled** API token
> - `iru_self_service_category` — requires the *Self Service* API scope
> - `iru_vulnerability`, `iru_vulnerability_detail` — requires the **Vulnerability Management** add-on (401 without the add-on)
>
> The following tables are confirmed working with real data: `iru_device`, `iru_device_detail`, `iru_device_app`, `iru_user`.

## Development

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
make install
```

Configure the plugin:

```sh
cp config/iru.spc ~/.steampipe/config/iru.spc
vi ~/.steampipe/config/iru.spc
```

## Testing

Run a smoke query against every table:

```sh
make test
```

The test script ([scripts/test_tables.sh](scripts/test_tables.sh)) builds the plugin, queries each table, and reports pass/fail/skip.

### Further reading

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Iru API reference](https://api-docs.kandji.io)
