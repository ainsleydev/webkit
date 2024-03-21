/**
 * Redirect to /admin when landing on the homepage for Strapi.
 *
 * @param _config
 * @param strapi
 */
export default (_config, { strapi }) => {
	const redirects = ['/', '/index.html'].map((path) => ({
		method: 'GET',
		path,
		handler: (ctx) => ctx.redirect('/admin'),
		config: { auth: false },
	}));

	strapi.server.routes(redirects);
};
