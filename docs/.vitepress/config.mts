import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
// Note: Mermaid support requires vitepress-plugin-mermaid (currently incompatible with VitePress 2.0 alpha)
// Diagrams are rendered as code blocks with mermaid class for client-side rendering
export default defineConfig({
    lang: 'en-GB',
    title: 'WebKit',
    description: 'Configuration-driven infrastructure for modern web projects',
    appearance: false,
    markdown: {
        theme: 'github-light'
    },
    sitemap: {
        hostname: 'https://webkit.ainsley.dev'
    },
    head: [
        ['link', { rel: 'icon', type: 'image/x-icon', href: '/favicon/favicon.ico' }],
        ['link', { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon/favicon-16x16.png' }],
        ['link', { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon/favicon-32x32.png' }],
        ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/favicon/apple-touch-icon.png' }],
        ['link', { rel: 'mask-icon', href: '/favicon/safari-pinned-tab.svg', color: '#70CBCE' }],
        ['link', { rel: 'manifest', href: '/favicon/site.webmanifest' }],
        ['meta', { name: 'msapplication-TileColor', content: '#70CBCE' }],
        ['meta', { name: 'msapplication-config', content: '/favicon/browserconfig.xml' }],
        ['meta', { name: 'theme-color', content: '#70CBCE' }],
    ],
    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        logo: '/logo-black.svg',
        nav: [
            { text: 'Home', link: '/' },
            { text: 'ainsley.dev', link: 'https://ainsley.dev' }
        ],
        sidebar: [
            {
                text: 'Getting started',
                items: [
                    { text: 'Installation', link: '/getting-started/installation' },
                    { text: 'Quick start', link: '/getting-started/quick-start' },
                    { text: 'Your first project', link: '/getting-started/your-first-project' },
                    { text: 'Core concepts', link: '/getting-started/core-concepts' }
                ]
            },
            {
                text: 'Manifest',
                items: [
                    { text: 'Overview', link: '/manifest/overview' },
                    { text: 'Project', link: '/manifest/project' },
                    { text: 'Apps', link: '/manifest/apps' },
                    { text: 'Resources', link: '/manifest/resources' },
                    { text: 'Monitoring', link: '/manifest/monitoring' },
                    { text: 'Environment variables', link: '/manifest/environment-variables' },
                    { text: 'Examples', link: '/manifest/examples' }
                ]
            },
            {
                text: 'CLI',
                items: [
                    { text: 'Command reference', link: '/cli/overview' },
                    { text: 'Validation', link: '/cli/validation' }
                ]
            },
            {
                text: 'Infrastructure',
                items: [
                    { text: 'Overview', link: '/infrastructure/overview' },
                    {
                        text: 'Providers',
                        collapsed: false,
                        items: [
                            { text: 'DigitalOcean', link: '/infrastructure/providers/digital-ocean' },
                            { text: 'Hetzner', link: '/infrastructure/providers/hetzner' },
                            { text: 'Backblaze B2', link: '/infrastructure/providers/backblaze-b2' },
                            { text: 'Turso', link: '/infrastructure/providers/turso' },
                            { text: 'Slack', link: '/infrastructure/providers/slack' }
                        ]
                    }
                ]
            }
        ],
        search: {
            provider: 'local'
        },
        editLink: {
            pattern: 'https://github.com/ainsleydev/webkit/edit/main/docs/:path',
            text: 'Edit this page on GitHub'
        },
        socialLinks: [
            { icon: 'github', link: 'https://github.com/ainsleydev/webkit' },
            { icon: 'instagram', link: 'https://www.instagram.com/ainsley.devltd/' },
            { icon: 'linkedin', link: 'https://www.linkedin.com/company/93587806' },
            { icon: 'twitter', link: 'https://x.com/ainsleydev/' }
        ],
        footer: {
            message: 'Released under the MIT Licence.',
            copyright: 'Copyright © 2024 ainsley.dev'
        }
    }
})
