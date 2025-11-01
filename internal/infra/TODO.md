# Infrastructure Import TODOs

## Hybrid Approach for Legacy Project Naming Conventions

### Problem Statement

The current import implementation assumes all resources follow webkit's naming convention:
- Full resource name: `${project_name}-${resource_name}` (e.g., "search-spares-db")
- Database prefix: lowercase with underscores (e.g., "search_spares_db")

However, legacy projects may have been created manually or with different naming conventions, making it impossible to import them without modification.

### Proposed Solution: Hybrid Import Strategy

Implement a flexible import system that tries webkit naming by default but provides overrides for legacy projects.

#### 1. Override Flags

Add granular flags to the `webkit infra import` command to override individual resource names:

```bash
webkit infra import \
  --resource=db \
  --id=cluster-123 \
  --db-user=custom_admin \
  --db-name=legacy_database \
  --db-pool=old_pool
```

**Flag Specifications:**

For PostgreSQL resources:
- `--db-user` - Override the database user name (default: `${project}_${resource}_admin`)
- `--db-name` - Override the database name (default: `${project}_${resource}`)
- `--db-pool` - Override the connection pool name (default: `${project}_${resource}_pool`)
- `--cluster-only` - Import only the cluster, skip user/database/pool (useful for partial imports)
- `--skip-firewall` - Skip firewall import even if `allowed_ips_addr` is configured

For S3/Spaces resources:
- `--bucket-name` - Override the bucket name (default: `${project}-${resource}`)

#### 2. Discovery Helper

Add a `--discover` flag that queries the cloud provider to show actual resource names before importing:

```bash
webkit infra import --resource=db --id=cluster-123 --discover
```

**Example Output:**
```
🔍 Discovering resources for cluster: cluster-123

DigitalOcean PostgreSQL Cluster
  Cluster ID: cluster-123

  Associated Resources:
  ✓ User:   legacy_db_admin
  ✓ Database: legacy_database
  ✓ Pool:   legacy_db_pool
  ✓ Firewall: configured (2 rules)

WebKit would expect:
  ✗ User:   search_spares_db_admin
  ✗ Database: search_spares_db
  ✗ Pool:   search_spares_db_pool

To import with current names, use:
  webkit infra import \
    --resource=db \
    --id=cluster-123 \
    --db-user=legacy_db_admin \
    --db-name=legacy_database \
    --db-pool=legacy_db_pool
```

#### 3. Auto-Detection (Future Enhancement)

Implement intelligent name detection that:
1. Tries webkit naming convention first
2. If import fails with "resource not found", automatically queries the provider API
3. Suggests the correct override flags based on discovered names
4. Asks user to confirm before retrying with discovered names

**Example Flow:**
```
⚠️  Import failed: database user "search_spares_db_admin" not found

🔍 Auto-detecting resource names...
   Found: legacy_db_admin

Would you like to retry with discovered names? (Y/n)
```

#### 4. Configuration File Support

Allow legacy project mappings to be stored in `app.json` to avoid repetitive flag usage:

```json
{
  "project": {
    "name": "search-spares"
  },
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "provider": "digitalocean",
      "config": {
        "allowed_ips_addr": ["185.16.161.205"]
      },
      "import_overrides": {
        "user_name": "legacy_db_admin",
        "database_name": "legacy_database",
        "pool_name": "legacy_db_pool"
      }
    }
  ]
}
```

### Implementation Considerations

1. **Backwards Compatibility**: Default behaviour must remain unchanged (webkit naming convention)

2. **Provider-Specific Logic**: Override flags and discovery should be provider-aware:
   - DigitalOcean: Use `doctl` for discovery
   - AWS: Use AWS SDK/CLI for discovery
   - B2: Use B2 SDK for discovery

3. **Validation**: When override flags are provided, validate them against the provider before attempting import

4. **Documentation**: Update CLI help text and docs to explain:
   - When to use override flags
   - How to discover existing resource names
   - Common legacy naming patterns

5. **Testing**: Add test cases for:
   - Import with override flags
   - Discovery helper output
   - Mixed webkit/legacy naming scenarios

### Priority

- **P0**: Override flags implementation (unblocks legacy project migrations)
- **P1**: Discovery helper (improves UX significantly)
- **P2**: Auto-detection (nice-to-have automation)
- **P3**: Configuration file support (convenience feature)

### Related Files

- `internal/infra/tf_import.go` - Core import address building
- `internal/infra/tf.go` - Import method implementation
- `internal/cmd/infra/import.go` - CLI command definition
- `internal/cmd/infra/import_test.go` - Import command tests
