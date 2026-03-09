# Table: iru_blueprint

List blueprints configured in your Iru tenant.

## Example queries

**List all blueprints:**

```sql
select id, name, active_devices, active_library_items, is_default_blueprint
from iru_blueprint
order by name;
```

**Find the default blueprint:**

```sql
select id, name, active_devices
from iru_blueprint
where is_default_blueprint = true;
```

**Blueprints with ADE enrollment enabled:**

```sql
select id, name, active_devices
from iru_blueprint
where is_dep_enrollment_enabled = true;
```
