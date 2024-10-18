import * as fs from 'node:fs';
import { configToJSONSchema, getPayload, type SanitizedConfig } from 'payload';
import { fieldMapper, type SchemaOptions, schemas } from '../plugin/schema.js';

/**
 * Creates JSON schema types of Payloads Collections & Globals
 */
export async function generateTypes(config: SanitizedConfig, opts: SchemaOptions): Promise<void> {
	console.log('Compiling JSON types for Collections and Globals...');

	const outputFile = (process.env.PAYLOAD_TS_OUTPUT_PATH || config.typescript.outputFile).replace(
		'.ts',
		'.json',
	);

	config.typescript = {
		...config.typescript,
		schema: schemas(opts)
	};

	// biome-ignore lint/style/noParameterAssign: Need to change field mapper.
	config = fieldMapper(config, opts);

	const payload = await getPayload({
		config,
		disableDBConnect: true,
		disableOnInit: true,
	});

	const jsonSchema = configToJSONSchema(payload.config, payload.db.defaultIDType);
	const prettyJSON = JSON.stringify(jsonSchema, null, 4);

	fs.writeFileSync(outputFile, prettyJSON);

	console.log(`JSON types written to: ${outputFile}`);

	delete process.env.GEN_GOLANG;
}
