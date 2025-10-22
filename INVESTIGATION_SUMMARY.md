# AGENTS.md Drift Investigation Summary

## Issue
Users are reporting drift warnings for AGENTS.md with the message "app.json changed, these files need regeneration" even when app.json hasn't been modified.

## Root Cause

### Template Change (Oct 22, 2025 @ 18:59)
Commit `0dae3d8` changed the AGENTS.md template:

```diff
-{{ .Content }}
+{{ .Content -}}
```

The `-` adds whitespace trimming, which removes a blank line when `.Content` is empty.

### Impact
1. **Old files** (generated before the change): Have an extra blank line
2. **New files** (generated after the change): Have one fewer blank line
3. **Drift detection**: Correctly identifies the difference as "outdated"
4. **Problem**: The message says "app.json changed" which is misleading - it's the template that changed, not app.json

## Why This Happens

The drift detection logic in `internal/manifest/drift.go:119-121`:

```go
if oldEntry.Hash == HashContent(actualData) {
    // File matches what we last generated, so app.json must have changed
    driftType = DriftReasonOutdated
}
```

This assumes that if a file needs regeneration but hasn't been manually modified, it must be because app.json changed. However, it could also be because:
- The template changed (this case)
- The WebKit CLI was updated with new template logic
- Default values or dependencies changed

## Verification

Comparing the template vs generated file:
- Template: 400 lines
- Generated: 399 lines (one line removed by `-` trimming)

The hash calculation is correct - no bugs in the hashing logic.

## Solution Options

### Option 1: Improve Error Message (Recommended)
Change the drift output message from:
```
app.json changed, these files need regeneration:
```

To:
```
Template or configuration changed, these files need regeneration:
```

### Option 2: Track Template Versions
Store the template version/hash in the manifest to distinguish between app.json changes and template changes.

### Option 3: Force Update
Add a note in the drift output suggesting users run `webkit update` to regenerate files with the latest templates.

## Recommendation

Implement **Option 1** immediately for better UX, and consider **Option 3** as an enhancement.

The drift detection is working correctly - this is expected behavior when templates change. Users just need to run `webkit update` once to regenerate AGENTS.md with the new template format.

## Files Involved

- `internal/templates/AGENTS.md` - The template that changed
- `internal/cmd/drift.go` - The drift detection command with misleading message
- `internal/manifest/drift.go` - The drift detection logic
- `internal/cmd/docs/agents.go` - The AGENTS.md generator
