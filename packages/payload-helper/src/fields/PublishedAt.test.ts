import type { DateField } from 'payload';
import { describe, expect, test, vi } from 'vitest';
import { PublishedAt } from './PublishedAt.js';

describe('PublishedAt', () => {
	test('returns a date field with correct defaults', () => {
		const field = PublishedAt();

		expect(field).toMatchObject({
			name: 'publishedAt',
			type: 'date',
			required: true,
			admin: {
				position: 'sidebar',
				date: {
					pickerAppearance: 'dayOnly',
				},
			},
		});
	});

	test('sets default value to current ISO date string', () => {
		const field = PublishedAt();
		const defaultValue = (field as DateField).defaultValue;

		expect(typeof defaultValue).toBe('function');

		const before = new Date().toISOString();
		const result = (defaultValue as () => string)();
		const after = new Date().toISOString();

		expect(result >= before).toBe(true);
		expect(result <= after).toBe(true);
	});

	test('sets date when status changes to published and value is empty', () => {
		const field = PublishedAt();
		const hooks = (field as DateField).hooks?.beforeChange;
		expect(hooks).toHaveLength(1);

		const before = new Date();
		const result = hooks?.[0]({
			siblingData: { _status: 'published' },
			value: undefined,
		} as unknown as Parameters<
			NonNullable<NonNullable<DateField['hooks']>['beforeChange']>[0]
		>[0]);
		const after = new Date();

		expect(result).toBeInstanceOf(Date);
		expect((result as Date).getTime()).toBeGreaterThanOrEqual(before.getTime());
		expect((result as Date).getTime()).toBeLessThanOrEqual(after.getTime());
	});

	test('preserves existing value when status is published', () => {
		const field = PublishedAt();
		const hooks = (field as DateField).hooks?.beforeChange;

		const existingDate = '2025-01-15T00:00:00.000Z';
		const result = hooks?.[0]({
			siblingData: { _status: 'published' },
			value: existingDate,
		} as unknown as Parameters<
			NonNullable<NonNullable<DateField['hooks']>['beforeChange']>[0]
		>[0]);

		expect(result).toBe(existingDate);
	});

	test('returns value as-is when status is not published', () => {
		const field = PublishedAt();
		const hooks = (field as DateField).hooks?.beforeChange;

		const result = hooks?.[0]({
			siblingData: { _status: 'draft' },
			value: undefined,
		} as unknown as Parameters<
			NonNullable<NonNullable<DateField['hooks']>['beforeChange']>[0]
		>[0]);

		expect(result).toBeUndefined();
	});

	test('applies overrides to the field', () => {
		const field = PublishedAt({
			overrides: {
				required: false,
				label: 'Date Published',
			},
		});

		expect(field).toMatchObject({
			name: 'publishedAt',
			type: 'date',
			required: false,
			label: 'Date Published',
		});
	});

	test('top-level overrides replace admin when admin is overridden', () => {
		const field = PublishedAt({
			overrides: {
				admin: {
					position: 'sidebar',
					description: 'Custom description',
				},
			},
		});

		expect((field as DateField).admin).toMatchObject({
			position: 'sidebar',
			description: 'Custom description',
		});
	});

	test('returns correct defaults with no arguments', () => {
		const field = PublishedAt();

		expect(field).toMatchObject({
			name: 'publishedAt',
			type: 'date',
			required: true,
		});
		expect((field as DateField).hooks?.beforeChange).toHaveLength(1);
		expect(typeof (field as DateField).defaultValue).toBe('function');
	});
});
