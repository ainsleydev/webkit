import { afterEach, beforeEach, describe, expect, test, vi } from 'vitest';

import { createLogger, createNoOpLogger } from './logger';

describe('createLogger', () => {
	const originalEnv = process.env;

	beforeEach(() => {
		vi.resetModules();
		process.env = { ...originalEnv };
	});

	afterEach(() => {
		process.env = originalEnv;
	});

	test('creates a logger with required service field', () => {
		const logger = createLogger({ service: 'test-service' });
		expect(logger).toBeDefined();
		expect(typeof logger.info).toBe('function');
		expect(typeof logger.error).toBe('function');
	});

	test('includes default company field', () => {
		const logger = createLogger({ service: 'test-service' });
		expect(logger.bindings()).toMatchObject({
			company: 'ainsley.dev',
			service: 'test-service',
		});
	});

	test('allows custom company field', () => {
		const logger = createLogger({
			service: 'test-service',
			company: 'custom-company',
		});
		expect(logger.bindings()).toMatchObject({
			company: 'custom-company',
		});
	});

	test('includes environment field from APP_ENV', () => {
		process.env.APP_ENV = 'staging';
		const logger = createLogger({ service: 'test-service' });
		expect(logger.bindings()).toMatchObject({
			environment: 'staging',
		});
	});

	test('falls back to NODE_ENV when APP_ENV not set', () => {
		delete process.env.APP_ENV;
		process.env.NODE_ENV = 'production';
		const logger = createLogger({ service: 'test-service' });
		expect(logger.bindings()).toMatchObject({
			environment: 'production',
		});
	});

	test('allows custom environment override', () => {
		process.env.APP_ENV = 'staging';
		const logger = createLogger({
			service: 'test-service',
			environment: 'custom-env',
		});
		expect(logger.bindings()).toMatchObject({
			environment: 'custom-env',
		});
	});

	test('uses LOG_LEVEL env var when set', () => {
		process.env.LOG_LEVEL = 'warn';
		const logger = createLogger({ service: 'test-service' });
		expect(logger.level).toBe('warn');
	});

	test('defaults to debug level in development', () => {
		delete process.env.LOG_LEVEL;
		process.env.NODE_ENV = 'development';
		const logger = createLogger({ service: 'test-service' });
		expect(logger.level).toBe('debug');
	});

	test('defaults to info level in production', () => {
		delete process.env.LOG_LEVEL;
		process.env.NODE_ENV = 'production';
		const logger = createLogger({ service: 'test-service' });
		expect(logger.level).toBe('info');
	});

	test('allows custom level override', () => {
		const logger = createLogger({
			service: 'test-service',
			level: 'error',
		});
		expect(logger.level).toBe('error');
	});

	test('supports child loggers with additional bindings', () => {
		const logger = createLogger({ service: 'test-service' });
		const child = logger.child({ request_id: 'abc-123' });
		expect(child.bindings()).toMatchObject({
			request_id: 'abc-123',
		});
	});
});

describe('createNoOpLogger', () => {
	test('creates a silent logger', () => {
		const logger = createNoOpLogger();
		expect(logger).toBeDefined();
		expect(logger.level).toBe('silent');
	});

	test('has all standard logging methods', () => {
		const logger = createNoOpLogger();
		expect(typeof logger.info).toBe('function');
		expect(typeof logger.error).toBe('function');
		expect(typeof logger.warn).toBe('function');
		expect(typeof logger.debug).toBe('function');
	});
});
