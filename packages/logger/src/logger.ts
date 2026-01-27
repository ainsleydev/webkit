import pino from 'pino';
import type { LogLevel, Logger, LoggerConfig } from './types.js';

const isBrowser = typeof window !== 'undefined';

const VALID_LEVELS: readonly LogLevel[] = [
	'trace',
	'debug',
	'info',
	'warn',
	'error',
	'fatal',
	'silent',
];

/**
 * Validates whether a string is a valid log level.
 * @param level - The level string to validate.
 * @returns True if the level is valid, false otherwise.
 */
function isValidLevel(level: string): level is LogLevel {
	return VALID_LEVELS.includes(level as LogLevel);
}

/**
 * Determines the default log level based on environment.
 * - Browser: always 'debug' (for devtools visibility).
 * - Node.js: LOG_LEVEL env var (if valid), or 'debug' in development, 'info' otherwise.
 * @param environment - The current environment name.
 * @returns The resolved log level.
 */
function resolveLevel(environment: string): LogLevel {
	if (isBrowser) {
		return 'debug';
	}
	const envLevel = process.env.LOG_LEVEL;
	if (envLevel && isValidLevel(envLevel)) {
		return envLevel;
	}
	const isDev = environment === 'development';
	return isDev ? 'debug' : 'info';
}

/**
 * Resolves the current environment from env vars.
 * @returns The environment name.
 */
function resolveEnvironment(): string {
	if (isBrowser) {
		return 'browser';
	}
	return process.env.APP_ENV || process.env.NODE_ENV || 'development';
}

/**
 * Creates a configured logger instance with structured JSON output.
 *
 * @param config - Logger configuration options.
 * @returns A configured pino logger instance.
 *
 * @example
 * ```typescript
 * import { createLogger } from '@ainsleydev/logger';
 *
 * const logger = createLogger({ service: 'my-app' });
 * logger.info('Server started', { attr: { port: 3000 } });
 * ```
 */
export function createLogger(config: LoggerConfig): Logger {
	const environment = config.environment ?? resolveEnvironment();
	const level = config.level ?? resolveLevel(environment);
	const company = config.company ?? 'ainsley.dev';

	const baseOptions: pino.LoggerOptions = {
		level,
		base: {
			company,
			service: config.service,
			environment,
		},
		timestamp: pino.stdTimeFunctions.isoTime,
		formatters: {
			level: (label) => ({ level: label }),
		},
	};

	if (isBrowser) {
		return pino({
			...baseOptions,
			browser: {
				asObject: true,
			},
		});
	}

	const isDev = environment === 'development';

	if (isDev) {
		let transport: pino.TransportSingleOptions | undefined;
		try {
			require.resolve('pino-pretty');
			transport = {
				target: 'pino-pretty',
				options: {
					colorize: true,
					translateTime: 'HH:MM:ss.l',
					ignore: 'pid,hostname',
				},
			};
		} catch {
			// pino-pretty not installed, use default JSON output.
		}

		if (transport) {
			return pino({
				...baseOptions,
				transport,
			});
		}
	}

	return pino(baseOptions);
}

/**
 * Creates a no-op logger that discards all output.
 * Useful for testing or when logging should be suppressed.
 *
 * @returns A logger that discards all output.
 */
export function createNoOpLogger(): Logger {
	return pino({ level: 'silent' });
}
