# Table: iru_prism_app

Cross-fleet application inventory from the Prism telemetry endpoint. Each row represents one application on one device, sourced from the real-time Prism data pipeline rather than the per-device apps endpoint.

> **Note:** Requires a Prism-enabled API token (`/api/v1/prism/apps`). Returns 401 with a standard token.

## Example queries

**Most widely installed applications across the fleet:**

```sql
select app_name, count(*) as device_count
from iru_prism_app
group by app_name
order by device_count desc
limit 20;
```

**All apps installed on a specific device:**

```sql
select app_name, bundle_id, version, path
from iru_prism_app
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by app_name;
```

**Devices running an outdated version of an app:**

```sql
select device_id, device_name, app_name, version
from iru_prism_app
where app_name = 'Zoom'
  and version < '6.0.0'
order by version;
```

**Recently updated app inventory entries:**

```sql
select device_name, app_name, version, updated_at
from iru_prism_app
order by updated_at desc
limit 50;
```
