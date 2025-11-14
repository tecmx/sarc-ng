/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check
import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  tutorialSidebar: [
    {type: "doc", id: "introduction"},
    {type: "doc", id: "getting-started"},
    {type: "doc", id: "architecture"},
    {type: "doc", id: "development"},
    {type: "doc", id: "deployment"},
  ],
  openApiSidebar: [
    {
      type: "category",
      label: "API Reference",
      link: {
        type: "generated-index",
        title: "API Reference",
        description:
          "Complete API reference for the SARC-NG (Resource Management and Scheduling System). This API provides endpoints for managing buildings, classrooms, resources, and scheduling through reservations and lessons.",
        slug: "/category/api-reference"
      },
      items: require("./content/api-reference/sidebar.js")
    }
  ]
};

export default sidebars;
