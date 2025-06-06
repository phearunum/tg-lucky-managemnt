import { defineConfig, loadEnv } from 'vite'
import path from 'path'
import createVitePlugins from './vite/plugins'

// https://vitejs.dev/config/
export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, process.cwd())

  const alias = {
    // set path
    '~': path.resolve(__dirname, './'),
    // set alias
    '@': path.resolve(__dirname, './src')
  }
  if (command === 'serve') {
    // resolve warning | esm-bundler build of vue-i18n.
    alias['vue-i18n'] = 'vue-i18n/dist/vue-i18n.cjs.js'
  }
  return {
    plugins: createVitePlugins(env, command === 'build'),
    resolve: {
      // https://cn.vitejs.dev/config/#resolve-alias
      alias: alias,
      // A list of extensions to omit when importing
      // https://cn.vitejs.dev/config/#resolve-extensions
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    css: {
      devSourcemap: true //Enabled in development mode
    },
    base: env.VITE_APP_ROUTER_PREFIX,
    // Packaging configuration
    build: {
      sourcemap: command === 'build' ? false : 'inline',
      outDir: 'dist', // build folder
      assetsDir: 'assets',
      rollupOptions: {
        output: {
          chunkFileNames: 'static/js/[name]-[hash].js',
          entryFileNames: 'static/js/[name]-[hash].js',
          assetFileNames: 'static/[ext]/[name]-[hash].[ext]'
        }
      }
    },
    // vite config
    server: {
      port: 8080,
      host: true,
      open: true,

      proxy: {
        // https://cn.vitejs.dev/config/#server-proxy
        '/dev-api': {
          target: env.VITE_APP_API_HOST,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/dev-api/, '')
        },
        '/msghub': {
          target: env.VITE_APP_API_HOST,
          ws: true,
          rewrite: (path) => path.replace(/^\/msgHub/, '')
        }
      }
    }
  }
})
