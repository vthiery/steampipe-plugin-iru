# Table: iru_prism_certificate

Cross-fleet certificate inventory from the Prism telemetry endpoint. Each row represents one certificate on one device.

> **Note:** Requires a Prism-enabled API token (`/api/v1/prism/certificates`). Returns 401 with a standard token.

## Example queries

**Expired certificates across the fleet:**

```sql
select device_id, device_name, common_name, issuer, expiry_date
from iru_prism_certificate
where is_expired = true
order by expiry_date;
```

**Certificates expiring within the next 30 days:**

```sql
select device_id, device_name, common_name, issuer, expiry_date
from iru_prism_certificate
where is_expired = false
  and expiry_date::date <= current_date + interval '30 days'
order by expiry_date;
```

**All certificates on a specific device:**

```sql
select common_name, issuer, subject, expiry_date, is_expired, key_usage
from iru_prism_certificate
where device_id = 'a6f07131-ed70-4b80-bed0-2c61f1d0b4df'
order by expiry_date;
```

**Count certificates per issuer:**

```sql
select issuer, count(*) as cert_count
from iru_prism_certificate
group by issuer
order by cert_count desc;
```
