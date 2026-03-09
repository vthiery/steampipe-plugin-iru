# Table: iru_library_item

List library items configured in your Iru tenant.

## Example queries

**List all active library items:**

```sql
select id, name, item_type, active
from iru_library_item
where active = true
order by item_type, name;
```

**Group library items by type:**

```sql
select item_type, count(*) as total
from iru_library_item
group by item_type
order by total desc;
```

**Find library items assigned to a specific blueprint:**

```sql
select id, name, item_type
from iru_library_item
where blueprint_ids @> '["YOUR_BLUEPRINT_ID"]'::jsonb;
```

**Get a specific library item by ID:**

```sql
select id, name, item_type, active, runs_on, updated_at
from iru_library_item
where id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df';
```
