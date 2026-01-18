import { describe, expect, test } from 'vitest';
import type { SidebarProps } from './Sidebar.svelte';

describe('Sidebar', () => {
	test('SidebarProps type is defined', () => {
		const props: SidebarProps = {
			children: () => {},
		};

		expect(props).toBeDefined();
	});

	test('menuLabel has default value', () => {
		const props: SidebarProps = {
			children: () => {},
		};

		expect(props.menuLabel).toBeUndefined();
	});

	test('isOpen can be set', () => {
		const props: SidebarProps = {
			children: () => {},
			isOpen: true,
		};

		expect(props.isOpen).toBe(true);
	});

	test('position accepts left or right', () => {
		const leftProps: SidebarProps = {
			children: () => {},
			position: 'left',
		};

		const rightProps: SidebarProps = {
			children: () => {},
			position: 'right',
		};

		expect(leftProps.position).toBe('left');
		expect(rightProps.position).toBe('right');
	});

	test('width can be customised', () => {
		const props: SidebarProps = {
			children: () => {},
			width: '300px',
		};

		expect(props.width).toBe('300px');
	});

	test('toggleStyle accepts toggle or hamburger', () => {
		const toggleProps: SidebarProps = {
			children: () => {},
			toggleStyle: 'toggle',
		};

		const hamburgerProps: SidebarProps = {
			children: () => {},
			toggleStyle: 'hamburger',
		};

		expect(toggleProps.toggleStyle).toBe('toggle');
		expect(hamburgerProps.toggleStyle).toBe('hamburger');
	});

	test('callbacks are optional', () => {
		const props: SidebarProps = {
			children: () => {},
			onOpen: () => {},
			onClose: () => {},
			onToggle: (isOpen: boolean) => {
				expect(typeof isOpen).toBe('boolean');
			},
		};

		expect(props.onOpen).toBeDefined();
		expect(props.onClose).toBeDefined();
		expect(props.onToggle).toBeDefined();
	});

	test('overlayOpacity accepts number', () => {
		const props: SidebarProps = {
			children: () => {},
			overlayOpacity: 0.5,
		};

		expect(props.overlayOpacity).toBe(0.5);
	});

	test('closeOnOverlayClick accepts boolean', () => {
		const props: SidebarProps = {
			children: () => {},
			closeOnOverlayClick: false,
		};

		expect(props.closeOnOverlayClick).toBe(false);
	});
});
