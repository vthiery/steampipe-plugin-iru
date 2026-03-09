# Table: iru_device_detail

Rich device details including MDM status, FileVault encryption, hardware overview, Kandji agent, and Activation Lock information.

> **Note:** This table fetches one API call per enrolled device. For large fleets it may take some time.

## Example queries

**FileVault coverage across all devices:**

```sql
select filevault_enabled, count(*) as device_count
from iru_device_detail
group by filevault_enabled
order by device_count desc;
```

**Devices with a specific Kandji agent version:**

```sql
select device_id, device_name, agent_version, agent_last_check_in
from iru_device_detail
where agent_version != '4.7.3';
```

**Devices without the personal recovery key escrowed:**

```sql
select device_id, device_name, filevault_enabled, filevault_recovery_key_type
from iru_device_detail
where filevault_enabled = true and filevault_prk_escrowed = false;
```

**Memory distribution across the fleet:**

```sql
select memory, count(*) as device_count
from iru_device_detail
group by memory
order by device_count desc;
```

**Get full details for a specific device:**

```sql
select *
from iru_device_detail
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df';
```
