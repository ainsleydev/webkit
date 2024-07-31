import path from 'node:path';
import { fileURLToPath } from 'node:url';
import type { Payload, PayloadRequest } from 'payload';
import { getFileByPath } from 'payload';
import type {Media, MediaSeed} from "./types.js";

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

export const uploadMedia = async (
	req: PayloadRequest,
	payload: Payload,
	media: MediaSeed,
): Promise<Media> => {
	try {
		const image = await getFileByPath(path.resolve(dirname, media.path));
		return await payload.create({
			collection: 'media',
			file: image,
			data: {
				alt: media.alt,
				//caption: media.caption,
			},
			req,
		}) as unknown as Media;
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
 */
export const up = async ({
	payload,
	req,
}: {
	payload: Payload;
	req: PayloadRequest;
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
};
