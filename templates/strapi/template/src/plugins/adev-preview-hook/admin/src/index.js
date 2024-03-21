//  import { StrapiApp } from '@strapi/admin/dist/admin/src/StrapiApp';
import { getFetchClient } from '@strapi/helper-plugin';

export default {
	/**
	 * Register the plugin with the Strapi application.
	 *
	 * @param {StrapiApp} app - The Strapi application instance.
	 */
	async register(app) {
		app.registerPlugin({
			id: 'adev-preview-hook',
			name: 'adev-preview-hook',
			apis: undefined,
			initializer: undefined,
			injectionZones: undefined,
			isReady: true,
		});
	},

	/**
	 * Set up the plugin's hooks and configurations during the bootstrap phase.
	 *
	 * @param {StrapiApp} app - The Strapi application instance.
	 */
	async bootstrap(app) {
		app.registerHook(
			'plugin/preview-button/before-build-url',
			async ({ data, draft, published }) => {
				const { get } = getFetchClient();
				let draftURL = draft.url;
				let publishedURL = draft.url;

				// Check if the draft URL needs to be populated with the 'channel.slug'
				// and the data exists within the fields.
				if (draft.url.includes('channel.slug') && data?.channel.length) {
					try {
						const response = await get(
							`/content-manager/collection-types/api::channel.channel/1`,
						);
						if (response.data.slug) {
							draftURL = draft.url.replace('{channel.slug}', response.data.slug);
							publishedURL = draft.url.replace('{channel.slug}', response.data.slug);
						}
					} catch (err) {
						console.log(err);
					}
				}

				return {
					draft: {
						url: draftURL,
						query: draft.query,
					},
					published: {
						url: publishedURL,
						query: published.query,
					},
				};
			},
		);
	},
};
