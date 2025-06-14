import { defineConfig, RsbuildConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginSass } from '@rsbuild/plugin-sass';
import { pluginTypedCSSModules } from '@rsbuild/plugin-typed-css-modules';
import { TanStackRouterRspack } from '@tanstack/router-plugin/rspack';
import path from 'path';

const config: RsbuildConfig = defineConfig({
    plugins: [pluginReact(), pluginSass(), pluginTypedCSSModules()],
    html: {
        template: path.resolve(__dirname, 'src/index.html'),
    },
    resolve: {
        alias: {
            'ui/*': './src/shared/components/ui/*',
            '@/*': './src/*',
        },
    },
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true,
                pathRewrite: { '^/api': '' },
            },
        },
    },
    tools: {
        rspack: {
            plugins: [
                TanStackRouterRspack({
                    target: 'react',
                    autoCodeSplitting: true,
                    routesDirectory: './src/app/routes',
                    generatedRouteTree: './src/app/routeTree.gen.ts',
                }),
            ],
        },
    },
});

export default config;
