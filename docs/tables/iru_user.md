# Table: iru_user

List users in your Iru tenant.

## Example queries

**List all active users:**

```sql
select id, name, email, active
from iru_user
where active = true
order by name;
```

**Find users with no devices assigned:**

```sql
select id, name, email
from iru_user
where device_count = 0;
```

**Get a specific user by ID:**

```sql
select id, name, email, device_count
from iru_user
where id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df';
```
