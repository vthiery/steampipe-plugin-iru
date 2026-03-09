# Table: iru_device_app

Lists all applications installed on each managed device. Each row represents one application on one device.

> **Note:** This table fetches one API call per enrolled device. For large fleets (~466 devices) you can expect 50 000+ rows.

## Example queries

**Most widely installed applications:**

```sql
select app_name, count(*) as device_count
from iru_device_app
group by app_name
order by device_count desc
limit 20;
```

**All apps on a specific device:**

```sql
select app_name, bundle_id, version, path
from iru_device_app
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by app_name;
```

**Devices running a vulnerable version of an app:**

```sql
select device_id, app_name, version, path
from iru_device_app
where app_name = 'Google Chrome'
  and version < '124.0.0.0';
```

**Apps not signed by Apple:**

```sql
select device_id, app_name, bundle_id, source, signature
from iru_device_app
where source != 'Apple'
order by device_id, app_name;
```

**Count distinct apps per device:**

```sql
select device_id, count(*) as app_count
from iru_device_app
group by device_id
order by app_count desc
limit 10;
```
