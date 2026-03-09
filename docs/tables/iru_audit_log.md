# Table: iru_audit_log

List audit log events from your Iru tenant.

## Example queries

**List the 20 most recent audit events:**

```sql
select id, event_type, actor_name, resource_type, resource_name, created_at
from iru_audit_log
order by created_at desc
limit 20;
```

**Count audit events by type:**

```sql
select event_type, count(*) as total
from iru_audit_log
group by event_type
order by total desc;
```

**Audit events performed by a specific actor:**

```sql
select id, event_type, action, resource_type, resource_name, created_at
from iru_audit_log
where actor_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by created_at desc;
```

**Device enrollment events:**

```sql
select id, actor_name, resource_name, ip_address, created_at
from iru_audit_log
where event_type = 'device.enrolled'
order by created_at desc;
```
