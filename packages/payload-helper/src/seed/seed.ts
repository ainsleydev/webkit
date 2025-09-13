import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import {
	type Payload,
	type PayloadRequest,
	type SanitizedConfig,
	commitTransaction,
	getPayload,
	initTransaction,
	killTransaction,
} from 'payload';
import env from '../util/env.js';

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

/**
 * A function that seeds the database with initial data.
 */
export type Seeder = (args: { payload: Payload; req: PayloadRequest }) => Promise<void>;

/**
 * Options for the seed function.
 * Note: You must use path.resolve for the paths, i.e. path.resolve(__dirname, 'path/to/file')
 */
export type SeedOptions = {
	config: SanitizedConfig;
	seeder: Seeder;
};

/**
 * Seeds the database with initial data.
 *
 * @param opts - The options for seeding.
 * @returns A promise that resolves when the seeding is complete.
 */
export const seed = (opts: {
	seeder: Seeder;
	config: SanitizedConfig;
}) => {
	const fn = async () => {
		process.env.PAYLOAD_DROP_DATABASE = 'true';

		const payload = await getPayload({
			config: opts.config,
		});
		const req = { payload } as PayloadRequest;

		await initTransaction(req);

		delete process.env.PAYLOAD_DROP_DATABASE;

		try {
			// Creating new tables
			payload.logger.info('Creating indexes...');
			try {
				if (payload.db.init) {
					await payload.db.init();
				}
			} catch (error) {
				payload.logger.error(`Creating database: ${error}`);
				return;
			}

			payload.logger.info('Migrating DB...');
			await payload.db.migrate();

			// Clearing local media
			if (!env.isProduction) {
				payload.logger.info('Clearing media...');
				const mediaDir = path.resolve(dirname, '../../media');
				if (fs.existsSync(mediaDir)) {
					fs.rmSync(mediaDir, { recursive: true, force: true });
				}
			}

			// Run user defined seed script
			await opts.seeder({ payload, req });

			await commitTransaction(req);

			payload.logger.info('Seed complete');
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Unknown error';
			payload.logger.error(`Seed failed: ${message}`);
			await killTransaction(req);
		}
	};

	fn()
		.then(() => process.exit(0))
		.catch((e) => {
			console.error(e);
			process.exit(1);
		});
};
