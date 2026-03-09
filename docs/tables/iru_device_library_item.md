# Table: iru_device_library_item

Lists the per-device library item assignment and compliance status for all managed devices.

> **Note:** Requires an API token with permission to read device library item status. Returns 403 with a standard read-only token.

## Example queries

**Library items with a failing status across the fleet:**

```sql
select device_id, name, type, status, reported_at
from iru_device_library_item
where status = 'fail'
order by reported_at desc;
```

**All library items for a specific device:**

```sql
select name, type, status, last_audit_run
from iru_device_library_item
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by name;
```

**Count failing items per device:**

```sql
select device_id, count(*) as fail_count
from iru_device_library_item
where status = 'fail'
group by device_id
order by fail_count desc;
```

**Library items of a specific type across the fleet:**

```sql
select device_id, name, status, reported_at
from iru_device_library_item
where type = 'profile'
order by device_id, name;
```
