# Backblaze B2

Backblaze B2 is an S3-compatible object storage service offering significantly lower pricing than most cloud providers. It's ideal for storing backups, media files, and large datasets.

## Authentication

Set your Backblaze B2 credentials as environment variables:

```bash
export B2_APPLICATION_KEY_ID="your-key-id"
export B2_APPLICATION_KEY="your-application-key"
```

Generate application keys at [secure.backblaze.com/app_keys.htm](https://secure.backblaze.com/app_keys.htm).

::: tip Application key scope
Create application keys with access limited to specific buckets for better security. WebKit only needs read/write access to the buckets it manages.
:::

## Object storage

B2 provides S3-compatible object storage at a fraction of typical cloud pricing.

### Configuration

```json
{
  "resources": [
    {
      "name": "backups",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPrivate"
      }
    }
  ]
}
```

### Bucket types

| Type | Description |
|------|-------------|
| `allPrivate` | All files are private (default) |
| `allPublic` | All files are publicly accessible |

### Pricing

B2's pricing is straightforward and significantly cheaper than alternatives:

| Resource | Price |
|----------|-------|
| Storage | $0.006/GB/month |
| Download | $0.01/GB (first 1GB free daily) |
| API calls | First 2,500 free, then $0.004/10,000 |

Compare to:
- AWS S3: ~$0.023/GB/month
- DigitalOcean Spaces: ~$0.020/GB/month
- Google Cloud Storage: ~$0.020/GB/month

### Outputs

| Output | Description |
|--------|-------------|
| `backups.bucket_name` | Bucket name |
| `backups.bucket_id` | Bucket ID |
| `backups.endpoint` | S3-compatible endpoint |

### Use in environment variables

```json
{
  "environment": {
    "production": {
      "BACKUP_BUCKET": {
        "source": "resource",
        "value": "backups.bucket_name"
      },
      "B2_ENDPOINT": {
        "source": "resource",
        "value": "backups.endpoint"
      }
    }
  }
}
```

## S3 compatibility

B2 provides an S3-compatible API, meaning you can use standard AWS SDKs and tools:

### Endpoint

The S3-compatible endpoint format is:

```
s3.{region}.backblazeb2.com
```

Regions:
- `us-west-004` - US West
- `eu-central-003` - EU Central

### SDK configuration

Configure your application to use B2's S3-compatible endpoint:

**Node.js (AWS SDK v3):**
```javascript
import { S3Client } from '@aws-sdk/client-s3';

const client = new S3Client({
  endpoint: 'https://s3.us-west-004.backblazeb2.com',
  region: 'us-west-004',
  credentials: {
    accessKeyId: process.env.B2_APPLICATION_KEY_ID,
    secretAccessKey: process.env.B2_APPLICATION_KEY,
  },
});
```

**Go:**
```go
cfg, _ := config.LoadDefaultConfig(ctx,
  config.WithRegion("us-west-004"),
  config.WithEndpointResolver(aws.EndpointResolverFunc(
    func(service, region string) (aws.Endpoint, error) {
      return aws.Endpoint{
        URL: "https://s3.us-west-004.backblazeb2.com",
      }, nil
    }),
  ),
)
```

## Use cases

### Backup storage

B2 is excellent for database and file backups:

```json
{
  "resources": [
    {
      "name": "db-backups",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPrivate",
        "lifecycle_rules": [
          {
            "prefix": "",
            "days_to_hide": 30,
            "days_to_delete": 90
          }
        ]
      }
    }
  ]
}
```

WebKit's backup workflows automatically upload to B2 when configured.

### Media storage

For user uploads and media files:

```json
{
  "resources": [
    {
      "name": "media",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPublic",
        "cors_rules": [
          {
            "allowed_origins": ["https://example.com"],
            "allowed_operations": ["s3_get", "s3_head"],
            "max_age_seconds": 3600
          }
        ]
      }
    }
  ]
}
```

### Archive storage

For long-term data retention:

```json
{
  "resources": [
    {
      "name": "archive",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPrivate",
        "lifecycle_rules": [
          {
            "prefix": "",
            "days_to_hide": 0,
            "days_to_delete": 365
          }
        ]
      }
    }
  ]
}
```

## Lifecycle rules

Automate file management with lifecycle rules:

```json
{
  "config": {
    "lifecycle_rules": [
      {
        "prefix": "logs/",
        "days_to_hide": 7,
        "days_to_delete": 30
      },
      {
        "prefix": "backups/",
        "days_to_hide": 30,
        "days_to_delete": 90
      }
    ]
  }
}
```

| Field | Description |
|-------|-------------|
| `prefix` | Apply rule to files with this prefix |
| `days_to_hide` | Days before hiding files (soft delete) |
| `days_to_delete` | Days before permanent deletion |

## Comparison with alternatives

| Feature | Backblaze B2 | DigitalOcean Spaces | AWS S3 |
|---------|--------------|---------------------|--------|
| Storage cost | $0.006/GB | $0.020/GB | $0.023/GB |
| Download cost | $0.01/GB | $0.01/GB | $0.09/GB |
| Free egress | 1GB/day | 1TB/month | None |
| CDN included | No | Yes | Extra cost |
| Regions | 2 | 8 | 20+ |
| S3 compatible | Yes | Yes | Native |

Choose B2 when:
- Cost is the primary concern
- You're storing backups or archives
- You don't need built-in CDN
- Download volume is moderate

## Example: Backup configuration

Complete backup setup with B2:

```json
{
  "resources": [
    {
      "name": "postgres",
      "type": "postgres",
      "provider": "digital_ocean",
      "config": {
        "size": "db-s-1vcpu-1gb"
      },
      "backup": {
        "enabled": true,
        "provider": "b2",
        "schedule": "0 3 * * *",
        "retention_days": 30
      }
    },
    {
      "name": "db-backups",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPrivate",
        "lifecycle_rules": [
          {
            "prefix": "",
            "days_to_hide": 30,
            "days_to_delete": 90
          }
        ]
      }
    }
  ]
}
```

## Further reading

- [B2 Cloud Storage documentation](https://www.backblaze.com/docs/cloud-storage)
- [S3-compatible API](https://www.backblaze.com/docs/cloud-storage-s3-compatible-api)
- [Pricing](https://www.backblaze.com/cloud-storage/pricing)
- [Application keys](https://www.backblaze.com/docs/cloud-storage-application-keys)
