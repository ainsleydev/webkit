import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import env from '../util/env.js';
import type { Payload, PayloadRequest } from 'payload';

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);


const postgresDrop = async (payload: Payload): Promise<void> => {
	// Drop all tables
	payload.logger.info('Dropping tables in database...');
	try {
		// @ts-expect-error
		const db = payload.db.pool;
		const client = await db.connect();
		await client.query('drop schema public cascade; create schema public;');
	} catch (error) {
		payload.logger.error(`Creating database: ${error}`);
		return;
	}
}

/**
 * Down script to remove all media and drop all tables.
 *
 * @param payload
 * @param req
 */
export const down = async ({
	payload,
	req,
}: {
	payload: Payload;
	req: PayloadRequest;
}): Promise<void> => {
	payload.logger.info('Running down script');

	// Clearing local media
	if (!env.isProduction) {
		payload.logger.info('Clearing media...');
		const mediaDir = path.resolve(dirname, '../../media');
		if (fs.existsSync(mediaDir)) {
			fs.rmSync(mediaDir, { recursive: true, force: true });
		}
	}

	await postgresDrop(payload);
};
