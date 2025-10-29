package gen

import _ "embed"

//go:embed docs/CODE_STYLE.md
var CodeStyle string

//go:embed docs/PAYLOAD.md
var Payload string

//go:embed docs/SVELTEKIT.md
var SvelteKit string
