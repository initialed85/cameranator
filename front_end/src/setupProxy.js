const {createProxyMiddleware} = require("http-proxy-middleware");

module.exports = function (app) {
    app.use(
        "/api/",
        createProxyMiddleware({
            target: "http://cameranator.chronos/",
            changeOrigin: true,
        })
    );

    app.use(
        "/events/",
        createProxyMiddleware({
            target: "http://cameranator.chronos/",
            changeOrigin: true,
        })
    );

    app.use(
        "/segments/",
        createProxyMiddleware({
            target: "http://cameranator.chronos/",
            changeOrigin: true,
        })
    );

    app.use(
        "/motion-stream/",
        createProxyMiddleware({
            target: "http://cameranator.chronos/",
            changeOrigin: true,
        })
    );
};
