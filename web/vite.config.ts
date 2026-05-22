import { fileURLToPath, URL } from 'node:url'
import { resolve, dirname } from 'node:path'
import { exec } from 'node:child_process'
import { promisify } from 'node:util'
const execPromised = promisify(exec)
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { defineConfig, type PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { PrimeVueResolver } from 'unplugin-vue-components/resolvers'
import { visualizer } from 'rollup-plugin-visualizer'
import checker from 'vite-plugin-checker'

import vueDevTools from 'vite-plugin-vue-devtools'
// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // const isProduction = true
  const isProduction = mode === 'production'

  process.env.VITE_APP_TITLE ??= isProduction
    ? 'Process model - Supply Chain Resilience'
    : 'Process Model - local'

  return {
    plugins: [
      vue(),
      Components({
        resolvers: [PrimeVueResolver()]
      }),
      vueDevTools(),
      visualizer() as PluginOption,
      VueI18nPlugin({
        /* options */
        // locale messages resource pre-compile option
        include: resolve(dirname(fileURLToPath(import.meta.url)), 'src/i18n/**')
      }),
      {
        name: 'translation', // the name of your custom plugin. Could be anything.
        handleHotUpdate: async ({ file }) => {
          if (file.includes('main_translation.js')) {
            await execPromised('npm run build-translation')
          }
        },
        buildStart: async () => {
          await execPromised('npm run build-translation')
        }
      },
      checker({
        typescript: true
      })
    ],
    build: {
      // sourcemap: true
      sourcemap: !isProduction
    },
    css: {
      devSourcemap: !isProduction
      // devSourcemap: true
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    }
  }
})
