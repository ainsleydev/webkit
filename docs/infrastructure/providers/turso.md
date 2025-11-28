# Turso

Turso provides edge-distributed SQLite databases, offering the simplicity of SQLite with the scalability of distributed systems. It's ideal for applications that benefit from low-latency reads across global locations.

## Authentication

Set your Turso credentials as environment variables:

```bash
export TURSO_API_TOKEN="your-api-token"
export TURSO_ORG_NAME="your-org-name"
```

Generate API tokens at [turso.tech/app](https://turso.tech/app) under your organisation settings.

## SQLite databases

Turso databases are libSQL (SQLite fork) databases replicated across edge locations.

### Configuration

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "group": "default"
      }
    }
  ]
}
```

### Groups

Turso organises databases into groups. Each group:
- Shares a primary location
- Has consistent replication settings
- Uses the same compute allocation

```json
{
  "config": {
    "group": "production",
    "primary_location": "lhr"
  }
}
```

### Locations

Turso has edge locations worldwide:

| Code | Location |
|------|----------|
| `ams` | Amsterdam |
| `arn` | Stockholm |
| `bog` | Bogota |
| `bom` | Mumbai |
| `cdg` | Paris |
| `den` | Denver |
| `dfw` | Dallas |
| `ewr` | Newark |
| `fra` | Frankfurt |
| `gru` | São Paulo |
| `hkg` | Hong Kong |
| `iad` | Washington DC |
| `jnb` | Johannesburg |
| `lax` | Los Angeles |
| `lhr` | London |
| `mad` | Madrid |
| `mia` | Miami |
| `nrt` | Tokyo |
| `ord` | Chicago |
| `otp` | Bucharest |
| `sea` | Seattle |
| `sin` | Singapore |
| `sjc` | San Jose |
| `syd` | Sydney |
| `waw` | Warsaw |

### Outputs

| Output | Description |
|--------|-------------|
| `db.url` | Database URL for connections |
| `db.auth_token` | Authentication token |
| `db.hostname` | Database hostname |

### Use in environment variables

```json
{
  "environment": {
    "production": {
      "DATABASE_URL": {
        "source": "resource",
        "value": "db.url"
      },
      "DATABASE_AUTH_TOKEN": {
        "source": "resource",
        "value": "db.auth_token"
      }
    }
  }
}
```

## Connecting to Turso

### Connection URL format

Turso URLs use the `libsql://` protocol:

```
libsql://your-database-your-org.turso.io
```

### SDK usage

**Node.js (@libsql/client):**
```javascript
import { createClient } from '@libsql/client';

const client = createClient({
  url: process.env.DATABASE_URL,
  authToken: process.env.DATABASE_AUTH_TOKEN,
});

const result = await client.execute('SELECT * FROM users');
```

**Go (libsql-client-go):**
```go
import "github.com/libsql/libsql-client-go/libsql"

connector, _ := libsql.NewConnector(
  os.Getenv("DATABASE_URL"),
  libsql.WithAuthToken(os.Getenv("DATABASE_AUTH_TOKEN")),
)
db := sql.OpenDB(connector)
```

**Rust (libsql):**
```rust
use libsql::Builder;

let db = Builder::new_remote(
    std::env::var("DATABASE_URL").unwrap(),
    std::env::var("DATABASE_AUTH_TOKEN").unwrap(),
)
.build()
.await?;
```

## Embedded replicas

Turso supports embedded replicas for ultra-low latency reads. The database is replicated to your application's memory.

```javascript
import { createClient } from '@libsql/client';

const client = createClient({
  url: 'file:local.db',
  syncUrl: process.env.DATABASE_URL,
  authToken: process.env.DATABASE_AUTH_TOKEN,
  syncInterval: 60, // Sync every 60 seconds
});

// Reads from local replica (microseconds)
const users = await client.execute('SELECT * FROM users');

// Writes go to primary, then sync
await client.execute('INSERT INTO users (name) VALUES (?)', ['Alice']);
await client.sync();
```

Configure embedded replicas in WebKit:

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "embedded_replicas": true
      }
    }
  ]
}
```

## Schema management

Turso supports schema-based multi-tenancy. Each schema is an isolated namespace within a database.

```json
{
  "config": {
    "allow_schema_attach": true
  }
}
```

### Creating schemas

```sql
-- Create a new schema for a tenant
ATTACH DATABASE 'tenant-123' AS tenant_123;

-- Create tables in the schema
CREATE TABLE tenant_123.users (
  id INTEGER PRIMARY KEY,
  name TEXT
);
```

### Querying schemas

```sql
-- Query specific tenant
SELECT * FROM tenant_123.users;

-- Join across schemas
SELECT * FROM tenant_a.orders
JOIN tenant_b.products ON ...;
```

## Pricing

Turso offers generous free tiers:

| Plan | Databases | Storage | Rows read | Rows written |
|------|-----------|---------|-----------|--------------|
| Starter (free) | 500 | 9GB | 1B/month | 25M/month |
| Scaler | Unlimited | 250GB | 100B/month | 100M/month |
| Enterprise | Unlimited | Custom | Custom | Custom |

## Use cases

### Edge-first applications

Deploy your database close to users:

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "primary_location": "lhr",
        "replicas": ["ams", "fra", "cdg"]
      }
    }
  ]
}
```

### Serverless applications

Turso's HTTP API works well with serverless:

```json
{
  "apps": [
    {
      "name": "api",
      "type": "svelte-kit",
      "environment": {
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "db.url"
          }
        }
      }
    }
  ],
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso"
    }
  ]
}
```

### Multi-tenant applications

Use schemas for tenant isolation:

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "allow_schema_attach": true,
        "enable_extensions": true
      }
    }
  ]
}
```

## Comparison with alternatives

| Feature | Turso | PlanetScale | Neon | Supabase |
|---------|-------|-------------|------|----------|
| Database | SQLite | MySQL | Postgres | Postgres |
| Edge replicas | Yes | No | No | No |
| Embedded mode | Yes | No | No | No |
| Free tier | Generous | Limited | Generous | Generous |
| Branching | No | Yes | Yes | No |
| Serverless | Yes | Yes | Yes | Yes |

Choose Turso when:
- You want SQLite simplicity
- Low-latency global reads matter
- You're building edge-first applications
- Embedded replicas would help your use case

## Example: SvelteKit with Turso

```json
{
  "project": {
    "name": "my-app",
    "title": "My App",
    "repo": "github.com/myorg/my-app"
  },
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app"
      },
      "environment": {
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "db.url"
          },
          "DATABASE_AUTH_TOKEN": {
            "source": "resource",
            "value": "db.auth_token"
          }
        }
      }
    }
  ],
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "group": "production",
        "primary_location": "lhr"
      }
    }
  ]
}
```

## Further reading

- [Turso documentation](https://docs.turso.tech/)
- [libSQL client libraries](https://docs.turso.tech/sdk)
- [Embedded replicas guide](https://docs.turso.tech/features/embedded-replicas)
- [Pricing](https://turso.tech/pricing)
