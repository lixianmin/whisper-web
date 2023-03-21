import {defineConfig} from 'vite'
import solidPlugin from 'vite-plugin-solid'
import path from 'path'

export default defineConfig({
    plugins: [solidPlugin()],
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