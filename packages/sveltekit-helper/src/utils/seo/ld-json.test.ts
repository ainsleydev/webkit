import { describe, expect, test } from 'vitest';
import { serializeSchema } from './ld-json';

describe('serializeSchema', () => {
	test('serialises a simple object to a JSON-LD script tag', () => {
		const input = { '@context': 'https://schema.org', '@type': 'WebSite', name: 'Test' };
		const got = serializeSchema(input);
		expect(got).toBe(
			'<script type="application/ld+json">{"@context":"https://schema.org","@type":"WebSite","name":"Test"}</script>',
		);
	});

	test('serialises a string value', () => {
		const got = serializeSchema('hello');
		expect(got).toBe('<script type="application/ld+json">"hello"</script>');
	});

	test('serialises null', () => {
		const got = serializeSchema(null);
		expect(got).toBe('<script type="application/ld+json">null</script>');
	});

	test('serialises nested objects', () => {
		const input = {
			'@context': 'https://schema.org',
			'@type': 'Organization',
			address: {
				'@type': 'PostalAddress',
				streetAddress: '123 Main St',
			},
		};
		const got = serializeSchema(input);
		expect(got).toContain('"@type":"Organization"');
		expect(got).toContain('"streetAddress":"123 Main St"');
		expect(got).toMatch(/^<script type="application\/ld\+json">.*<\/script>$/);
	});
});
