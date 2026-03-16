---
"@ainsleydev/sveltekit-helper": patch
---

Fix `TableOfContents` component erroring at runtime due to nested HTML comments in the `@component` doc block closing the outer comment early, causing Svelte to parse example markup as real template code.
