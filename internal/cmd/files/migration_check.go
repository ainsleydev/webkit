package files

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
)

// MigrationCheckScript scaffolds a dependency check script for Payload CMS apps.
// This script ensures dependencies are up-to-date before running migrations.
func MigrationCheckScript(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter to only Payload apps.
	payloadApps := appDef.GetAppsByType(appdef.AppTypePayload)

	for _, app := range payloadApps {
		scriptPath := filepath.Join(app.Path, "scripts", "check-deps.js")

		// Scaffold the check script.
		err := input.Generator().Bytes(scriptPath, []byte(checkDepsScript),
			scaffold.WithTracking(manifest.SourceApp(app.Name)),
			scaffold.WithScaffoldMode(),
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("creating migration check script for app %s", app.Name))
		}
	}

	return nil
}

// checkDepsScript is a simple Node.js script that checks if dependencies are in sync.
// It compares the modification time of pnpm-lock.yaml against node_modules.
const checkDepsScript = `const fs = require('fs');
const path = require('path');

try {
	const lockFile = path.join(__dirname, '..', 'pnpm-lock.yaml');
	const nodeModules = path.join(__dirname, '..', 'node_modules');

	if (!fs.existsSync(lockFile)) {
		console.error('❌ pnpm-lock.yaml not found');
		process.exit(1);
	}

	if (!fs.existsSync(nodeModules)) {
		console.error('❌ node_modules not found. Run: pnpm install');
		process.exit(1);
	}

	const lockStat = fs.statSync(lockFile);
	const nodeModulesStat = fs.statSync(nodeModules);

	if (lockStat.mtimeMs > nodeModulesStat.mtimeMs) {
		console.error('❌ Dependencies out of sync!');
		console.error('   pnpm-lock.yaml is newer than node_modules');
		console.error('   Run: pnpm install');
		process.exit(1);
	}

	console.log('✅ Dependencies are in sync');
} catch (err) {
	console.error('❌ Error checking dependencies:', err.message);
	process.exit(1);
}
`
