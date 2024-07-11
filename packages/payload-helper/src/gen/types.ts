import fs from 'fs';
import path from 'path';
import payload from 'payload';
import getLogger from 'payload/dist/utilities/logger';
import { configToJSONSchema } from './schema';

/**
 * Creates JSON schema types of Payloads Collections & Globals
 *
 * Should probably fork this:
 * https://github.com/payloadcms/payload/blob/b700208b98e0b49ef86c7cbfa18751110da4e7b3/packages/payload/src/utilities/configToJSONSchema.ts#L545
 */
export async function generateTypes(): Promise<void> {
	const logger = getLogger();

	logger.info('Compiling JSON types for Collections and Globals...');

	await payload.init({
		config: payload.config,
		disableDBConnect: true,
		disableOnInit: true,
		local: true,
		secret: '--unused--',
	});

	const jsonSchema = configToJSONSchema(payload.config, payload.db.defaultIDType, payload);
	const prettyJSON = JSON.stringify(jsonSchema, null, 4);
	const outFile = './types/payload.json';

	fs.writeFileSync(path.resolve(__dirname, outFile), prettyJSON);

	logger.info('JSON types written to: ' + outFile);
}

generateTypes().catch((e) => console.log(e));
