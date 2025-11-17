import { describe, expect, test } from 'vitest';
import { z } from 'zod';

import { generateFormSchema } from './schema';

describe('generateFormSchema', () => {
	test('returns empty object schema when no fields provided', () => {
		const schema = generateFormSchema(null);
		expect(schema).toBeInstanceOf(z.ZodObject);
		expect(schema.shape).toEqual({});
	});

	test('returns empty object schema for empty array', () => {
		const schema = generateFormSchema([]);
		expect(schema).toBeInstanceOf(z.ZodObject);
		expect(schema.shape).toEqual({});
	});

	test('generates required text field schema', () => {
		const fields = [
			{
				blockType: 'text' as const,
				name: 'firstName',
				label: 'First Name',
				required: true,
			},
		];

		const schema = generateFormSchema(fields);
		const result = schema.safeParse({ firstName: '' });
		expect(result.success).toBe(false);

		const validResult = schema.safeParse({ firstName: 'John' });
		expect(validResult.success).toBe(true);
	});

	test('generates optional text field schema', () => {
		const fields = [
			{
				blockType: 'text' as const,
				name: 'middleName',
				label: 'Middle Name',
				required: false,
			},
		];

		const schema = generateFormSchema(fields);
		const result = schema.safeParse({});
		expect(result.success).toBe(true);
	});

	test('generates required email field schema', () => {
		const fields = [
			{
				blockType: 'email' as const,
				name: 'email',
				label: 'Email',
				required: true,
			},
		];

		const schema = generateFormSchema(fields);
		const invalidResult = schema.safeParse({ email: 'not-an-email' });
		expect(invalidResult.success).toBe(false);

		const validResult = schema.safeParse({ email: 'test@example.com' });
		expect(validResult.success).toBe(true);
	});

	test('generates required checkbox field schema', () => {
		const fields = [
			{
				blockType: 'checkbox' as const,
				name: 'terms',
				label: 'Accept Terms',
				required: true,
			},
		];

		const schema = generateFormSchema(fields);
		const falseResult = schema.safeParse({ terms: false });
		expect(falseResult.success).toBe(false);

		const trueResult = schema.safeParse({ terms: true });
		expect(trueResult.success).toBe(true);
	});

	test('generates schema with multiple fields', () => {
		const fields = [
			{
				blockType: 'text' as const,
				name: 'name',
				label: 'Name',
				required: true,
			},
			{
				blockType: 'email' as const,
				name: 'email',
				label: 'Email',
				required: true,
			},
			{
				blockType: 'textarea' as const,
				name: 'message',
				label: 'Message',
				required: false,
			},
		];

		const schema = generateFormSchema(fields);
		const result = schema.safeParse({
			name: 'John',
			email: 'john@example.com',
		});
		expect(result.success).toBe(true);
	});
});
