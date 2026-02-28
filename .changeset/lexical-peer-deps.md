---
"@ainsleydev/payload-helper": patch
---

Moving lexical and @lexical/* to peerDependencies with a >=0.35.0 range so consuming projects inherit whatever version their Payload install resolves, avoiding version conflicts across different Payload releases.
