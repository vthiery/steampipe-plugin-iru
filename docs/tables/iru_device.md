# Table: iru_device

List and inspect devices enrolled in your Iru tenant.

## Example queries

**List all enrolled Mac devices:**

```sql
select device_id, device_name, serial_number, os_version, last_check_in
from iru_device
where platform = 'Mac'
order by last_check_in desc;
```

**Find devices not enrolled in MDM:**

```sql
select device_id, device_name, serial_number, platform
from iru_device
where mdm_enabled = false;
```

**Devices grouped by blueprint:**

```sql
select blueprint_name, count(*) as device_count
from iru_device
group by blueprint_name
order by device_count desc;
```

**Find missing devices:**

```sql
select device_id, device_name, serial_number, last_check_in
from iru_device
where is_missing = true;
```

**Get a specific device by ID:**

```sql
select device_id, device_name, serial_number, os_version, assigned_user_email
from iru_device
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df';
```
