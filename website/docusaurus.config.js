// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'KNDP',
  tagline: 'KNDP is cool',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://github.com/',
  baseUrl: '/',    
  organizationName: 'web-seven', // Usually your GitHub org/user name.
  projectName: 'kndp', // Usually your repo name.
  deploymentBranch: "gh-pages",
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  
   

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  

};

module.exports=config;