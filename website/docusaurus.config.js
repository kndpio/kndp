module.exports = {
  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: false,
    },
    tableOfContents: {
      minHeadingLevel: 2,
      maxHeadingLevel: 5,
    },

  },
  
};


const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');


/** @type {import('@docusaurus/types').Config} */
const config = {
  title: ' ',
  tagline: '',
  favicon: 'img/',

  // Set the production url of your site here
  url: 'https://kndp.io/',
  baseUrl: '/',
  organizationName: 'web-seven', // Usually your GitHub org/user name.
  projectName: 'kndp', // Usually your repo name.
  deploymentBranch: 'gh-pages',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
        },
        blog: {
          showReadingTime: true,
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],


  themeConfig: {
    liveCodeBlock: {
      playgroundPosition: 'bottom',
    },
    docs: {
      sidebar: {
        hideable: false,
        autoCollapseCategories: false,
      },
    },
    colorMode: {
      defaultMode: 'light',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    prism: {
      additionalLanguages: ['java', 'latex'],
      magicComments: [
        {
          className: 'theme-code-block-highlighted-line',
          line: 'highlight-next-line',
          block: { start: 'highlight-start', end: 'highlight-end' },
        },
      ],
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
      
    },
    algolia: {
      appId: ' ',
      apiKey: ' ',
      indexName: 'docusaurus-2',
    }, 

    navbar: {
      hideOnScroll: true,
      title: 'KNDP',
      logo: {
        alt: '',
        src: '/img/logo-draft.svg',
      },
      items: [
        { to: '/why-kndp', label: 'Why KNDP?', position: 'right' },
        {
          type: 'doc',
          position: 'right',
          docId: 'overview',
          label: 'Documentation',
        },         
        { to: '/community', label: 'Community', position: 'right' },
        { to: '/blog', label: 'Blog', position: 'right' },
        {
          href: 'https://github.com/web-seven/kndp',
          position: 'right',
          className: 'header-github-link',
          'aria-label': 'GitHub repository',
        },
      ].filter(Boolean),
    },
    footer: {
      logo: {
        alt: 'Web Seven Logo',
        src: 'img/logo_text.png',
        href: 'https://web7.dev',
        width: 70,
      },
      copyright: `Copyright © ${new Date().getFullYear()} Web Seven SRL. Built with Docusaurus.`,
    },
  },
};

module.exports = config;
