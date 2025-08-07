import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';
import * as path from 'path';

const config: Config = {
  title: 'SARC-NG',
  tagline: 'Resource Reservation and Management API',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://sarc-ng.example.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'tecmx', // Usually your GitHub org/user name.
  projectName: 'sarc-ng', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: path.resolve(__dirname, './sidebars.ts'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/tecmx/sarc-ng/tree/main/docs/',
          routeBasePath: '/', // Serve the docs at the site's root
        },
        blog: false, // Disable the blog feature
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  plugins: [
    [
      require.resolve('docusaurus-plugin-openapi-docs'),
      {
        id: "api",
        docsPluginId: "classic",
        config: {
          sarc: {
            specPath: "static/openapi.yaml",
            outputDir: "docs/api-reference",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "tag",
            },
          },
        },
      },
    ],
    [
      '@docusaurus/plugin-content-docs',
      {
        id: 'api-docs',
        path: 'docs/api-reference',
        routeBasePath: 'api-reference',
        sidebarPath: require.resolve('./docs/api-reference/sidebar.ts'),
      },
    ],
  ],

  themes: ['docusaurus-theme-openapi-docs'],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/sarc-ng-social-card.jpg',
    navbar: {
      title: 'SARC-NG',
      logo: {
        alt: 'SARC-NG Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Documentation',
        },
        {
          to: '/api-reference',
          position: 'left',
          label: 'API Reference',
        },
        {
          href: 'https://github.com/tecmx/sarc-ng',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Introduction',
              to: '/',
            },
            {
              label: 'API Reference',
              to: '/api-reference',
            },
          ],
        },
        {
          title: 'Resources',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/tecmx/sarc-ng',
            },
            {
              label: 'Issues',
              href: 'https://github.com/tecmx/sarc-ng/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} SARC-NG Project. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'go', 'json', 'yaml'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
