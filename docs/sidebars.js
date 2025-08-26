/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  tutorialSidebar: [
    'introduction',
    {
      type: 'category',
      label: 'Getting Started',
      items: [
        'getting-started/overview',
        'getting-started/installation',
        'getting-started/quick-start',
        'getting-started/deployment-modes',
      ],
    },
    {
      type: 'category',
      label: 'Core Concepts',
      items: [
        'concepts/architecture',
        'concepts/accounts',
        'concepts/transactions',
        'concepts/tokens',
        'concepts/idempotency',
        'concepts/security',
      ],
    },
    {
      type: 'category',
      label: 'Deployment Guides',
      items: [
        'deployment/lightweight-mode',
        'deployment/production-setup',
        'deployment/network-configuration',
        'deployment/security-best-practices',
      ],
    },
    {
      type: 'category',
      label: 'API Reference',
      items: [
        'api-reference/overview',
        'api-reference/accounts',
        'api-reference/transactions',
        'api-reference/tokens',
        'api-reference/system',
      ],
    },
    {
      type: 'category',
      label: 'Advanced Topics',
      items: [
        'advanced/key-management',
        'advanced/worker-pools',
        'advanced/chain-events',
        'advanced/troubleshooting',
      ],
    },
    {
      type: 'category',
      label: 'Examples',
      items: [
        'examples/basic-usage',
        'examples/token-transfers',
        'examples/integration-patterns',
      ],
    },
  ],
};

export default sidebars;