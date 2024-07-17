import { PayloadHandler, PayloadRequest } from 'payload';

/**
 * Find a document in the given collection by its slug.
 *
 * @param collection
 */
export const findBySlug = (collection: string): PayloadHandler => {
	return async (req: PayloadRequest): Promise<Response> => {
		try {
			const data = await req.payload.find({
				collection,
				where: {
					slug: {
						equals: req?.routeParams?.slug ?? '',
					},
				},
				limit: 1,
			});
			if (data.docs.length === 0) {
				return new Response(JSON.stringify({ error: 'not found' }), { status: 404 });
			} else {
				return new Response(JSON.stringify(data.docs[0]), { status: 200 });
			}
		} catch (error) {
			console.error('Error occurred while fetching document:', error);
			return new Response(JSON.stringify({ error: 'Internal server error' }), {
				status: 500,
			});
		}
	};
};
