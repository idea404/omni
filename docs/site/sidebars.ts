import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // tutorialSidebar: [{type: 'autogenerated', dirName: '.', exclude: ['omni', 'home']}],
  aboutSidebar: [
    "home",
    {
      type: "category",
      label: "About Omni",
      collapsible: false,
      items: [
        {
          type: "category",
          label: "Background",
          items: [
            {
              type: "autogenerated",
              dirName: "background",
            },
          ],
        },
        {
          type: "autogenerated",
          dirName: "omni",
        },
      ],
    },
    // {
    //   type: "category",
    //   label: "Protocol",
    //   items: [
    //     {
    //       type: "autogenerated",
    //       dirName: "protocol",
    //     }
    //   ]
    // },
    // {
    //   type: "category",
    //   label: "Developers",
    //   collapsible: false,
    //   items: [
    //     {
    //       type: "autogenerated",
    //       dirName: "developers",
    //     }
    //   ]
    // },
    {
      type: "doc",
      id: "resources/glossary",
    },
  ],
  protocolSidebar: [
    {
      type: "autogenerated",
      dirName: "protocol",
    },
  ],
  developersSidebar: [
    {
      type: "autogenerated",
      dirName: "developers",
    },
  ],
  operatorsSidebar: [
    {
      type: "autogenerated",
      dirName: "operators",
    },
  ],
};

export default sidebars;