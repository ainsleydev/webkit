import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
	lang: 'en-GB',
	title: "WebKit",
	description: "Infrastructure as code framework for full-stack web applications",
	sitemap: {
		hostname: 'https://webkit.ainsley.dev'
	},
	markdown: {
		theme: {
			light: 'github-light',
			dark: 'github-dark'
		}
	},
	themeConfig: {
		// https://vitepress.dev/reference/default-theme-config
		nav: [
			{text: 'Home', link: '/'},
			{text: 'Getting started', link: '/getting-started/installation'},
			{text: 'Docs', link: '/core-concepts/overview'},
			{text: 'CLI', link: '/cli/overview'},
		],
		sidebar: [
			{
				text: 'Getting started',
				collapsed: false,
				items: [
					{text: 'Installation', link: '/getting-started/installation'},
					{text: 'Quick start', link: '/getting-started/quick-start'},
					{text: 'Your first project', link: '/getting-started/your-first-project'}
				]
			},
			{
				text: 'Core concepts',
				collapsed: false,
				items: [
					{text: 'Overview', link: '/core-concepts/overview'}
				]
			},
			{
				text: 'Manifest',
				collapsed: false,
				items: [
					{text: 'Overview', link: '/manifest/overview'},
					{text: 'Project', link: '/manifest/project'},
					{text: 'Apps', link: '/manifest/apps'},
					{text: 'Resources', link: '/manifest/resources'},
					{text: 'Environment variables', link: '/manifest/environment-variables'},
					{text: 'Examples', link: '/manifest/examples'}
				]
			},
			{
				text: 'CLI reference',
				collapsed: true,
				items: [
					{text: 'Overview', link: '/cli/overview'},
					{text: 'webkit update', link: '/cli/webkit-update'}
				]
			},
			{
				text: 'Infrastructure',
				collapsed: true,
				items: [
					{text: 'Overview', link: '/infrastructure/overview'}
				]
			}
		],
		search: {
			provider: 'local'
		},
		socialLinks: [
			{icon: 'github', link: 'https://github.com/ainsleydev/webkit'},
		],
		footer: {
			message: 'Released under the BSD-3 Clause License',
			copyright: 'Copyright © 2023-present ainsley.dev'
		},
		editLink: {
			pattern: 'https://github.com/ainsleydev/webkit/edit/feat-major/app-definitions/docs/:path',
			text: 'Edit this page on GitHub'
		}
	}
})
