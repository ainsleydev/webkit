import * as fs from 'node:fs';
import { configToJSONSchema, getPayload } from 'payload';
import { findConfig, importConfig } from 'payload/node';

/**
 * Creates JSON schema types of Payloads Collections & Globals
 */
export async function generateTypes(outFile: string): Promise<void> {
	console.log('Compiling JSON types for Collections and Globals...');

	let configPath = '';
	try {
		configPath = findConfig();
	} catch (e) {
		console.log('Error finding config: ' + e);
		return;
	}

	// Set the environment variable to generate Golang types.
	process.env.GEN_GOLANG = 'true';

	const config = await importConfig(configPath);
	const outputFile = (process.env.PAYLOAD_TS_OUTPUT_PATH || config.typescript.outputFile).replace(
		'.ts',
		'.json',
	);

	const payload = await getPayload({
		config,
		disableDBConnect: true,
		disableOnInit: true,
		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
		// @ts-ignore
		local: true,
		secret: '--unused--',
	});

	const jsonSchema = configToJSONSchema(payload.config, payload.db.defaultIDType);
	const prettyJSON = JSON.stringify(jsonSchema, null, 4);

	fs.writeFileSync(outputFile, prettyJSON);

	console.log(`JSON types written to: ${outputFile}`);

	delete process.env.GEN_GOLANG;
}
