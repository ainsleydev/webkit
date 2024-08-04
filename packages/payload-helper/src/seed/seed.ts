import dotenv from 'dotenv';
import {
	type PayloadRequest,
	commitTransaction,
	getPayload,
	initTransaction,
	killTransaction,
	type Payload,
} from 'payload';
import { importConfig } from 'payload/node';
import { down } from './down.js';
import { up } from './up.js';

/**
 * A function that seeds the database with initial data.
 */
export type Seeder = (args: { payload: Payload; req: PayloadRequest }) => Promise<void>;

/**
 * Options for the seed function.
 * Note: You must use path.resolve for the paths, i.e. path.resolve(__dirname, 'path/to/file')
 */
export type SeedOptions = {
	envPath: string;
	configPath: string;
	dbAdapter: DBAdapter;
	seeder: Seeder;
};

/**
 * The database adapter to use, which will remove and recreate the database.
 */
export enum DBAdapter {
	Postgres = 'postgres',
}

/**
 * Seeds the database with initial data.
 *
 * @param opts - The options for seeding.
 * @returns A promise that resolves when the seeding is complete.
 */
export const seed = (opts: SeedOptions) => {
	const fn = async () => {
		dotenv.config({
			path: opts.envPath,
		});

		for (const fn of [down, up]) {
			if (fn === down) {
				process.env.PAYLOAD_DROP_DATABASE = 'true';
			} else {
				delete process.env.PAYLOAD_DROP_DATABASE; // Ensure it is not set for other functions
			}

			const config = await importConfig(opts.configPath);
			const payload = await getPayload({ config });
			const req = { payload } as PayloadRequest;

			await initTransaction(req);

			try {
				await fn({ payload, req, seeder: opts.seeder });
				payload.logger.info('Seed complete');
				await commitTransaction(req);
			} catch (err) {
				const message = err instanceof Error ? err.message : 'Unknown error';
				payload.logger.error(`Seed failed: ${message}`);
				await killTransaction(req);
			}
		}
	};

	fn().then(() => process.exit(0)).catch((e) => {
		console.error(e);
		process.exit(1);
	});
};
