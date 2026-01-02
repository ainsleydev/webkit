// Dependency sync check for Payload CMS migrations.
// Compares pnpm-lock.yaml against the cached lockfile in node_modules
// to prevent CI failures from outdated dependencies.

const fs = require('fs');
const path = require('path');

try {
	const lockFile = path.join(__dirname, '..', 'pnpm-lock.yaml');
	const nodeModulesLock = path.join(__dirname, '..', 'node_modules', '.pnpm', 'lock.yaml');

	if (!fs.existsSync(lockFile)) {
		console.error('❌ pnpm-lock.yaml not found');
		process.exit(1);
	}

	if (!fs.existsSync(nodeModulesLock)) {
		console.error('❌ Dependencies not installed. Run: pnpm install');
		process.exit(1);
	}

	const lockContent = fs.readFileSync(lockFile, 'utf8');
	const nodeModulesContent = fs.readFileSync(nodeModulesLock, 'utf8');

	// pnpm creates identical lockfiles, so strict equality check is correct.
	// This catches actual dependency mismatches without false positives from timestamps.
	if (lockContent !== nodeModulesContent) {
		console.error('❌ Dependencies out of sync!');
		console.error('   Run: pnpm install');
		process.exit(1);
	}

	console.log('✅ Dependencies are in sync');
} catch (err) {
	console.error('❌ Error checking dependencies:', err.message);
	process.exit(1);
}
