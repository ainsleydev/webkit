import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
    lang: 'en-GB',
    title: "WebKit",
    description: "ainsley.dev's App Manifest",
    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
            {text: 'Home', link: '/'},
            {text: 'ainsley.dev', link: '/'},
        ],
        sidebar: [
            {
                text: 'Intro',
                items: [
                    {text: 'Summary', link: '/intro.md'}
                ]
            },
            {
                text: 'Manifest',
                items: [
                    {text: 'Overview', link: '/manifest/overview.md'},
                    {text: 'Project', link: '/manifest/project.md'},
                    {text: 'Resources', link: '/manifest/resources.md'},
                    {text: 'Apps', link: '/manifest/apps.md'},
                    {text: 'Environment Variables', link: '/manifest/environment-variables.md'}
                ]
            }
        ],
        socialLinks: [
            {icon: 'github', link: 'https://github.com/ainsleydev'},
            {icon: 'instagram', link: 'https://www.instagram.com/ainsley.devltd/'},
            {icon: 'linkedin', link: 'https://www.linkedin.com/company/93587806'},
            {icon: 'twitter', link: 'https://x.com/ainsleydev/'},
            {icon: 'facebook', link: 'https://www.facebook.com/ainsley.devltd/'},
        ]
    }
})
