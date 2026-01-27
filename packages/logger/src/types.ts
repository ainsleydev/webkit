import type { Logger as PinoLogger } from 'pino';

/**
 * Log levels supported by the logger.
 */
export type LogLevel = 'trace' | 'debug' | 'info' | 'warn' | 'error' | 'fatal' | 'silent';

/**
 * Configuration options for creating a logger instance.
 */
export interface LoggerConfig {
	/** The name of the service or application. */
	service: string;
	/** The company name. Defaults to "ainsley.dev". */
	company?: string;
	/** The environment (e.g., "development", "production"). Defaults to APP_ENV or NODE_ENV. */
	environment?: string;
	/** The log level. Defaults based on environment. */
	level?: LogLevel;
}

/**
 * The logger instance type.
 */
export type Logger = PinoLogger;
