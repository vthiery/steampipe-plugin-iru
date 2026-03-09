# Table: iru_device_parameter

Lists per-device MDM parameter compliance status for all managed devices. Each row represents one MDM parameter check on one device.

> **Note:** Requires an API token with permission to read device status. Returns 403 with a standard read-only token.

## Example queries

**Parameters with a failing status across the fleet:**

```sql
select device_id, category, name, status
from iru_device_parameter
where status = 'fail'
order by device_id, category, name;
```

**All parameters for a specific device:**

```sql
select category, subcategory, name, status
from iru_device_parameter
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by category, name;
```

**Count failing parameters per device:**

```sql
select device_id, count(*) as fail_count
from iru_device_parameter
where status = 'fail'
group by device_id
order by fail_count desc
limit 20;
```

**Compliance summary by category:**

```sql
select category, status, count(*) as count
from iru_device_parameter
group by category, status
order by category, status;
```
