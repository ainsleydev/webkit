# @ainsleydev/sveltekit-helper

## 0.4.0

### Minor Changes

- 69e11e6: Add PayloadSEO and PayloadFooter components for rendering head meta tags (Open Graph, Twitter Card, canonical URL, JSON-LD structured data) and footer code injection from Payload CMS settings. Extract `serializeSchema` as a generic utility under `utils/seo`. Add `resolveItems` helper for resolving Payload relationship fields.

## 0.3.2

### Patch Changes

- bac8c57: Fix icon colour CSS custom properties (`--_alert-icon-colour`, `--_notice-icon-colour`) not being applied to icons in Alert and Notice components. Add `hideIcon` prop to both components to optionally hide the icon.

## 0.3.1

### Patch Changes

- 41f24c1: Fixing CSS variables

## 0.3.0

### Minor Changes

- ad8a9ab: Adding notification components

## 0.2.1

### Patch Changes

- 45fc51f: Upgrade svelte-hamburgers to v5.0.0 for Svelte 5 compatibility and enable runes mode

## 0.2.0

### Minor Changes

- 2e83635: Adding Sidebar and Hamburger components for mobile-first navigation. Includes customisable props, CSS variable support with inline fallbacks, and toggle/hamburger display modes. New dependency: svelte-hamburgers

## 0.1.5

### Patch Changes

- 99c58f8: Improve CSS variable override flexibility with fallback pattern. Adds mobile-specific variables (--row-gap-mobile, --col-gap-mobile) to allow responsive customization without media query conflicts. Includes comprehensive CSS specificity documentation

## 0.1.4

### Patch Changes

- 9e57cdb: Fixing container example

## 0.1.3

### Patch Changes

- 0a894ab: Replaced $$restProps with Svelte 5 runes mode rest destructuring using $props()

## 0.1.2

### Patch Changes

- d0c2bd2: Fixed Svelte components to use Svelte 5 runes mode instead of legacy export let syntax

## 0.1.1

### Patch Changes

- e91e490: Fix build process to include Svelte component files in distribution
