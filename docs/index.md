# Iru Plugin for Steampipe

Use SQL to query devices, users, blueprints, threats, and more from [Iru](https://iru.com) (formerly Kandji), the Apple device management platform.

- **[Get started →](https://github.com/vthiery/steampipe-plugin-iru)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-iru/issues)

## Installation

Clone and build the plugin:

```sh
git clone https://github.com/vthiery/steampipe-plugin-iru.git
cd steampipe-plugin-iru
make install
```

## Configuration

Copy the sample config:

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
  # You can find this in Settings > Access in the Iru web app.
  url = "https://acme.api.kandji.io"

  # API token generated in Settings > Access > API Token.
  token = "YOUR_API_TOKEN"
}
```

## Tables

### Device inventory

| Table | Description |
|-------|-------------|
| [iru_device](tables/iru_device.md) | List and inspect devices enrolled in your Iru tenant. |
| [iru_device_detail](tables/iru_device_detail.md) | Rich device details: MDM, FileVault, hardware, Kandji agent, and Activation Lock. |
| [iru_device_app](tables/iru_device_app.md) | Applications installed on each managed device. |
| [iru_device_library_item](tables/iru_device_library_item.md) | Per-device library item assignment and compliance status. |
| [iru_device_parameter](tables/iru_device_parameter.md) | Per-device MDM parameter compliance status. |

### Users, blueprints & tags

| Table | Description |
|-------|-------------|
| [iru_user](tables/iru_user.md) | List users in your Iru tenant. |
| [iru_blueprint](tables/iru_blueprint.md) | List blueprints configured in your Iru tenant. |
| [iru_tag](tables/iru_tag.md) | List tags configured in your Iru tenant. |

### Library, threats & audit

| Table | Description |
|-------|-------------|
| [iru_library_item](tables/iru_library_item.md) | List library items configured in your Iru tenant. |
| [iru_threat](tables/iru_threat.md) | List threat events detected on devices in your Iru tenant. |
| [iru_audit_log](tables/iru_audit_log.md) | List audit log events from your Iru tenant. |

### Prism telemetry

> These tables require a Prism-enabled API token (see Settings → Access in the Iru web app).

| Table | Description |
|-------|-------------|
| [iru_prism_app](tables/iru_prism_app.md) | Cross-fleet application inventory from the Prism pipeline. |
| [iru_prism_certificate](tables/iru_prism_certificate.md) | Cross-fleet certificate inventory from the Prism pipeline. |
| [iru_prism_filevault](tables/iru_prism_filevault.md) | Cross-fleet FileVault encryption status from the Prism pipeline. |

### Self Service

| Table | Description |
|-------|-------------|
| [iru_self_service_category](tables/iru_self_service_category.md) | Categories defined in the Iru Self Service catalog. |

### Vulnerability Management

> These tables require the **Vulnerability Management add-on** to be enabled on your Iru tenant.

| Table | Description |
|-------|-------------|
| [iru_vulnerability](tables/iru_vulnerability.md) | CVE-grouped vulnerabilities detected across the tenant. |
| [iru_vulnerability_detail](tables/iru_vulnerability_detail.md) | Full CVE detail (description, EPSS score, publish dates) for a specific CVE ID. |
