import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      '@lib': path.resolve(__dirname, 'src/lib'),
      '@go': path.resolve(__dirname, 'wailsjs/go'),
      '@runtime': path.resolve(__dirname, 'wailsjs/runtime'),
    },
  },
})
