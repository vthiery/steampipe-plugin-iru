# Table: iru_self_service_category

Lists categories defined in the Iru Self Service catalog.

> **Note:** Requires an API token with Self Service read permissions. Returns 403 with a standard read-only token.

## Example queries

**List all Self Service categories:**

```sql
select id, name, description, sort_order
from iru_self_service_category
order by sort_order;
```

**Get a specific category by ID:**

```sql
select id, name, description
from iru_self_service_category
where id = 'abc12345-0000-0000-0000-000000000000';
```

**Categories with no description:**

```sql
select id, name, sort_order
from iru_self_service_category
where description = '' or description is null
order by sort_order;
```
