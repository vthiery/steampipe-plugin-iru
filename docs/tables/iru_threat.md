# Table: iru_threat

List threat events detected on devices in your Iru tenant.

## Example queries

**List all active (unresolved) threats:**

```sql
select threat_id, device_name, file_name, classification, status, date_time
from iru_threat
where status != 'resolved'
order by date_time desc;
```

**Count threats by classification:**

```sql
select classification, count(*) as total
from iru_threat
group by classification
order by total desc;
```

**Threats on a specific device:**

```sql
select threat_id, file_name, sha256, classification, status, date_time
from iru_threat
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by date_time desc;
```

**List blocked malware by family:**

```sql
select malware_family, count(*) as total
from iru_threat
where status = 'blocked'
  and malware_family != ''
group by malware_family
order by total desc;
```
