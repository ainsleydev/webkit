# @ainsleydev/logger

Structured JSON logging for Node.js and browser environments using [pino](https://github.com/pinojs/pino).

## Installation

```bash
pnpm add @ainsleydev/logger
```

## Usage

```typescript
import { createLogger } from '@ainsleydev/logger';

const logger = createLogger({ service: 'my-app' });

logger.info('Server started', { port: 3000 });
logger.error('Connection failed', { error: err.message });
```

### With request ID

```typescript
const child = logger.child({ request_id: 'abc-123' });
child.info('Handling request');
```

## Configuration

```typescript
const logger = createLogger({
  service: 'my-app',           // Required
  company: 'ainsley.dev',      // Default: "ainsley.dev"
  environment: 'production',   // Default: APP_ENV || NODE_ENV || "development"
  level: 'info',               // Default: see below
});
```

### Log levels

- **Browser**: Always `debug`.
- **Node.js**: `LOG_LEVEL` env var, or `debug` in development, `info` in production.

Valid levels: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `silent`.

## Output

JSON output to stdout:

```json
{"time":"2025-01-27T10:30:00.000Z","level":"info","msg":"Server started","company":"ainsley.dev","service":"my-app","environment":"production","port":3000}
```

### Pretty printing (development)

Install `pino-pretty` for human-readable output in development:

```bash
pnpm add -D pino-pretty
```
