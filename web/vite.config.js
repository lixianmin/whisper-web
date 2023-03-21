import {defineConfig} from 'vite'
import solidPlugin from 'vite-plugin-solid'
import path from 'path'

export default defineConfig({
    plugins: [solidPlugin(),
        {
            name: 'configure-response-headers',
            configureServer: server => {
                server.middlewares.use((_req, res, next) => {
                    res.setHeader('Cross-Origin-Embedder-Policy', 'require-corp');
                    res.setHeader('Cross-Origin-Opener-Policy', 'same-origin');
                    next();
                });
            }
        }
    ],
    server: {
        port: 3000,
        https: {
            key: path.resolve(__dirname, '../res/ssl/localhost.key'),
            cert: path.resolve(__dirname, '../res/ssl/localhost.crt'),
        }
    },
    build: {
        target: 'esnext',
    },
});