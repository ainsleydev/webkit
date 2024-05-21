import type {GlobalConfig} from 'payload/types'
import {languages} from "./locales";
import {TelephoneField} from '@nouance/payload-better-fields-plugin'
import {validateURL} from "../util/validation";

export const Settings: GlobalConfig = {
    slug: 'settings',
    typescript: {
        interface: 'Settings',
    },
    graphQL: {
        name: 'Settings',
    },
    access: {
        read: () => true,
    },
    fields: [
        {
            type: 'tabs',
            tabs: [
                {
                    label: 'Global',
                    description: 'Configure global settings for the website.',
                    fields: [
                        {
                            type: 'row',
                            fields: [
                                {
                                    name: 'siteName',
                                    type: 'text',
                                    label: 'Site Name',
                                    admin: {
                                        width: '50%',
                                        description: 'Add a site name for the website, this will be outputted in the Open Graph schema as well as a suffix for the meta title.',
                                    }
                                },
                                {
                                    name: 'locale',
                                    type: 'select',
                                    label: 'Locale',
                                    defaultValue: 'en_GB',
                                    options: languages.map(l => {
                                        return {
                                            label: l.name,
                                            value: l.code,
                                        }
                                    }),
                                    admin: {
                                        width: '50%',
                                        description: 'Add a locale for the website, this will be outputted in the Open Graph schema and the top level HTML tag. Defaults to en_GB.',
                                    }
                                },
                            ],
                        },
                        {
                            name: 'tagLine',
                            type: 'textarea',
                            label: 'Tag Line',
                            admin: {
                                description: 'In a few words, explain what this site is about',
                            }
                        },
                        {
                            name: 'logo',
                            type: 'upload',
                            relationTo: 'media',
                            filterOptions: {
                                mimeType: {
                                    contains: 'image',
                                },
                            },
                            admin: {
                                description: 'Add a logo for the website that will be displayed in the header & across the website.',
                            },
                        },
                        {
                            name: 'robots',
                            type: 'textarea',
                            label: 'Robots.txt',
                            admin: {
                                description: 'Robots.txt is a text file webmasters create to instruct web robots (typically search engine robots) how to crawl pages on their website.',
                            }
                        }
                    ]
                },
                {
                    label: 'Code Injection',
                    description: 'Code injection allows you to inject a small snippet of HTML into your site. It can be a css override, analytics of a block javascript.',
                    fields: [
                        {
                            name: 'codeInjection',
                            type: 'group',
                            interfaceName: 'CodeInjection',
                            fields: [
                                {
                                    name: 'head',
                                    type: 'code',
                                    label: 'Head',
                                    admin: {
                                        language: 'html',
                                        description: 'Outputs code within the <head> of the website.',
                                    }
                                },
                                {
                                    name: 'footer',
                                    type: 'code',
                                    label: 'Footer',
                                    admin: {
                                        language: 'html',
                                        description: 'Outputs code in the footer of the website.',
                                    }
                                },
                            ]
                        }
                    ]
                },
                {
                    label: 'Contact Details',
                    fields: [
                        {
                           name: 'contact',
                            type: 'group',
                            interfaceName: 'Contact',
                            fields: [
                                {
                                    type: 'row',
                                    fields: [
                                        {
                                            name: 'email',
                                            type: 'email',
                                            label: 'Email',
                                            admin: {
                                                width: '40%',
                                            },
                                        },
                                        ...TelephoneField({
                                            name: 'telephone',
                                            admin: {
                                                width: '50%',
                                                placeholder: '+44 1732 123456',
                                            },
                                        }),
                                    ]
                                },
                                {
                                    name: 'address',
                                    type: 'textarea',
                                    label: 'Address',
                                },
                                {
                                    type: 'row',
                                    fields: [
                                        {
                                            name: 'linkedIn',
                                            type: 'text',
                                            label: 'LinkedIn',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a LinkedIn URL for the website.',
                                            },
                                        },
                                        {
                                            name: 'x',
                                            type: 'text',
                                            label: 'X',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a X (Twitter) URL for the website.',
                                            },
                                        },
                                        {
                                            name: 'facebook',
                                            type: 'text',
                                            label: 'Facebook',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a Facebook URL for the website.',
                                            },
                                        },
                                        {
                                            name: 'instagram',
                                            type: 'text',
                                            label: 'Instagram',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a Instagram URL for the website.',
                                            },
                                        },
                                        {
                                            name: 'youtube',
                                            type: 'text',
                                            label: 'Youtube',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a Youtube URL for the website.',
                                            },
                                        },
                                        {
                                            name: 'tiktok',
                                            type: 'text',
                                            label: 'TikTok',
                                            validate: validateURL,
                                            admin: {
                                                width: '50%',
                                                description: 'Add a TikTok URL for the website.',
                                            },
                                        },
                                    ]
                                }
                            ]
                        }
                    ],
                },
                {
                    label: 'Maintenance',
                    fields: [
                        {
                            name: 'maintenance',
                            type: 'group',
                            interfaceName: 'Maintenance',
                            fields: [
                                {
                                    name: 'enabled',
                                    type: 'checkbox',
                                    label: 'Enable',
                                    admin: {
                                        description: 'Enable maintenance mode for the site, this will use a maintenance page template and not include any of the sites functioanlity.',
                                    },
                                },
                                {
                                    name: 'title',
                                    type: 'text',
                                    label: 'Title',
                                    admin: {
                                        description: 'Add a title for the maintenance page.',
                                    }
                                },
                                {
                                    name: 'content',
                                    type: 'textarea',
                                    label: 'Content',
                                    admin: {
                                        description: 'Add content for the maintenance page, it will appear beneath the title.',
                                    }
                                }
                            ]
                        }
                    ]
                }
            ],
        }
    ],
}
