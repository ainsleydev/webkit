import { describe, expect, test } from 'vitest';
import type { HamburgerProps } from './Hamburger.svelte';

describe('Hamburger', () => {
	test('HamburgerProps type is defined', () => {
		const props: HamburgerProps = {};

		expect(props).toBeDefined();
	});

	test('isOpen can be set', () => {
		const props: HamburgerProps = {
			isOpen: true,
		};

		expect(props.isOpen).toBe(true);
	});

	test('gap can be customised', () => {
		const props: HamburgerProps = {
			gap: '1rem',
		};

		expect(props.gap).toBe('1rem');
	});

	test('ariaLabel can be set', () => {
		const props: HamburgerProps = {
			ariaLabel: 'Open navigation',
		};

		expect(props.ariaLabel).toBe('Open navigation');
	});

	test('onChange callback is optional', () => {
		const props: HamburgerProps = {
			onChange: (isOpen: boolean) => {
				expect(typeof isOpen).toBe('boolean');
			},
		};

		expect(props.onChange).toBeDefined();
	});

	test('class prop can be set', () => {
		const props: HamburgerProps = {
			class: 'custom-class',
		};

		expect(props.class).toBe('custom-class');
	});
});
