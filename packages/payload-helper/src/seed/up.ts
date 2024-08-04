import path from 'node:path';
import type { Payload, PayloadRequest } from 'payload';
import { getFileByPath } from 'payload';
import { htmlToLexical } from '../util/lexical.js';
import type { Seeder } from './seed.js';
import type { Media, MediaSeed } from './types.js';

/**
 *
 * @param req
 * @param payload
 * @param dirname
 * @param media
 */
export const uploadMedia = async (
	req: PayloadRequest,
	payload: Payload,
	dirname: string,
	media: MediaSeed,
): Promise<Media> => {
	try {
		const image = await getFileByPath(path.resolve(dirname, media.path));
		const caption = media.caption ? await htmlToLexical(media.caption) : null;

		return (await payload.create({
			collection: 'media',
			file: image,
			data: {
				alt: media.alt,
				caption: caption,
			},
			req,
		})) as unknown as Media;
	} catch (error) {
		payload.logger.error(`Uploading media: ${error}`);
		throw error;
	}
};

/**
 * Up script to create tables and seed data.
 *
 * @param payload
 * @param req
 * @param seeder
 */
export const up = async ({
	payload,
	req,
	seeder,
}: {
	payload: Payload;
	req: PayloadRequest;
	seeder: Seeder;
}): Promise<void> => {
	payload.logger.info('Running up script');

	await payload.init({
		config: payload.config,
	});

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

	await seeder({ payload, req });
};
