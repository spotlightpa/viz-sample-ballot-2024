const pluginTailwind = require("prettier-plugin-tailwindcss");
const pluginGoTemplate = require("prettier-plugin-go-template");

module.exports = {
  plugins: [pluginGoTemplate, pluginTailwind],
  pluginSearchDirs: false,
  tailwindConfig: "./tailwind.config.js",
  overrides: [
    {
      files: ["*.html"],
      options: {
        parser: "go-template",
      },
    },
  ],
};
