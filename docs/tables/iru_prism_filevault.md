# Table: iru_prism_filevault

Cross-fleet FileVault encryption status from the Prism telemetry endpoint. Provides richer detail than `iru_device_detail` (e.g. key rotation scheduling) with cursor-based pagination.

> **Note:** Requires a Prism-enabled API token (`/api/v1/prism/filevault`). Returns 401 with a standard token.

## Example queries

**Devices without FileVault enabled:**

```sql
select device_id, device_name, recovery_key_type, updated_at
from iru_prism_filevault
where filevault_enabled = false
order by device_name;
```

**Devices where the personal recovery key is not escrowed:**

```sql
select device_id, device_name, filevault_enabled, recovery_key_type
from iru_prism_filevault
where filevault_enabled = true and prk_escrowed = false
order by device_name;
```

**FileVault coverage summary:**

```sql
select filevault_enabled, count(*) as device_count
from iru_prism_filevault
group by filevault_enabled;
```

**Devices with key regeneration required:**

```sql
select device_id, device_name, recovery_key_type, next_rotation
from iru_prism_filevault
where regen_required = true
order by device_name;
```
