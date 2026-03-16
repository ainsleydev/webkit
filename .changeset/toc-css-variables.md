---
"@ainsleydev/sveltekit-helper": patch
---

feat(toc): expose CSS variables on TableOfContents for active colour, border colour and offset

- `--toc-colour-active` — overrides active/hover link colour (fallback: `--token-text-action`)
- `--toc-border-colour` — overrides border colour (fallback: `--colour-light-600`)
- `--toc-border-offset` — overrides `margin-left` and `padding-left` on the border variant (fallback: `$size-48`)
