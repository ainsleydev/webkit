module.exports = {
	/**
	 * Sends an automated email when a form type is
	 * inserted into the database.
	 *
	 * @param result
	 */
	async afterCreate({ result }) {
		const optsPage = await strapi.entityService.findPage('api::option.option');
		if (!optsPage.results.length) {
			throw new Error('Could not retrieve option data from form lifecycle.');
		}

		const options = optsPage.results[0];
		const tpl = `
<h3>New email form submission on ${options.siteName}</h3>
<br/>
<p><b>Name:</b> ${result.name ?? ''}</p>
<p><b>Email:</b> ${result.email ?? ''}</p>
<p><b>Message:</b> ${result.message ?? ''}</p>
<p><b>Privacy:</b> ${result.privacy ?? false}</p>`;

		await strapi
			.plugin('email')
			.service('email')
			.send({
				to: options.email,
				from: process.env.EMAIL_FROM ?? 'hello@audits.com',
				subject: 'New form submission on ' + (options.siteName ?? 'your website'),
				text: tpl.replace(/<[^>]*>?/gm, ''),
				html: tpl,
			});
	},
};
