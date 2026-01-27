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
 * Base fields included in every log entry.
 */
export interface LogFields {
	/** The company name. */
	company: string;
	/** The service or application name. */
	service: string;
	/** The environment name. */
	environment: string;
	/** Optional request identifier. */
	request_id?: string;
}

/**
 * Additional attributes to include in log entries.
 * These are grouped under the "attr" key in JSON output.
 */
export type LogAttributes = Record<string, unknown>;

/**
 * The logger instance type.
 */
export type Logger = PinoLogger;
