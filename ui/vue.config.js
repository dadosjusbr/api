var path = require("path");
var PrerenderSpaPlugin = require("prerender-spa-plugin");

module.exports = {
  publicPath: "/",
  productionSourceMap: false,
  configureWebpack: {
    plugins: [
      new PrerenderSpaPlugin({
        // Absolute path to compiled SPA
        staticDir: path.resolve(__dirname, "./dist"),
        // List of routes to prerender
        routes: ["/", "/sobre", "/orgao/mppb/2020/1"],
      }),
    ],
  },
};
