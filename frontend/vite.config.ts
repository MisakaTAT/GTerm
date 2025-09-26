import { resolve } from 'node:path';
import vue from '@vitejs/plugin-vue';
import AutoImport from 'unplugin-auto-import/vite';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      imports: ['vue', 'vue-router'],
      dirs: [],
      resolvers: [NaiveUiResolver()],
      vueTemplate: true,
      dts: 'src/auto-imports.d.ts',
      eslintrc: {
        enabled: false,
        filepath: resolve(__dirname, 'eslint.config.autoImport.json'),
        globalsPropValue: 'readonly',
      },
    }),
    Components({
      dts: 'src/components.d.ts',
      deep: true,
      dirs: ['src/components'],
      extensions: ['vue', 'tsx'],
      resolvers: [NaiveUiResolver()],
    }),
  ],
  server: {
    hmr: {
      host: 'localhost',
      protocol: 'ws',
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@wailsApp': resolve(__dirname, 'bindings'),
    },
  },
});
