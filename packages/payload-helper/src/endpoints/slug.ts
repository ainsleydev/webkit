import { PayloadHandler } from 'payload/config';

/**
 * Find a document in the given collection by its slug.
 *
 * @param collection
 */
export const findBySlug = (collection: string): PayloadHandler => {
	return async (req, res) => {
		try {
			const data = await req.payload.find({
				collection,
				where: {
					slug: {
						equals: req.params.slug,
					},
				},
				limit: 1,
			});
			if (data.docs.length === 0) {
				res.status(404).send({ error: 'not found' });
			} else {
				res.status(200).send(data.docs[0]);
			}
		} catch (error) {
			console.error('Error occurred while fetching document:', error);
			res.status(500).send({ error: 'Internal server error' });
		}
	};
};
