const esbuild = require("esbuild")
const { sassPlugin } = require("esbuild-sass-plugin")

// Function to perform the build
async function buildProject(entryPoints, outdir, watch = false) {
    const buildOptions = {
        entryPoints: entryPoints,
        bundle: true,
        external: ["/static/*", "desandro-matches-selector", "ev-emitter", "get-size", "fizzy-ui-utils", "./item", "outlayer", "popper.js"],
        outdir: outdir,
        sourcemap: watch,
        loader: {
            ".js": "jsx",
            ".scss": "css",
        },
        minify: true,
        plugins: [
            sassPlugin({type: "css"}),
        ],
        target: ["esnext"],
        define: {},
    }

    if (watch) {
        const context = await esbuild.context(buildOptions)
        context.watch()
        console.log("Watching for changes...")
    } else {
        await esbuild.build(buildOptions)
    }
}
const isWatchMode = process.argv.includes("--watch")
const entryPoints = [
    "js/index.ts",
    "js/workers/service.worker.ts",
    "css/main.scss",
    "css/*.css",
]
buildProject(entryPoints, "dist", isWatchMode)
