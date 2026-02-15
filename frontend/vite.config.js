import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import fs from 'fs'

// Read version from versioninfo.json
let appVersion = '0.0.0'
try {
  const versionInfoPath = resolve(__dirname, '../versioninfo.json')
  if (fs.existsSync(versionInfoPath)) {
    const data = fs.readFileSync(versionInfoPath, 'utf-8')
    const json = JSON.parse(data)
    appVersion = json.StringFileInfo?.ProductVersion || '0.0.0'
  }
} catch (e) {
  console.warn('Failed to read version info:', e)
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  base: './',
  define: {
    __APP_VERSION__: JSON.stringify(appVersion)
  },
  build: {
    outDir: '../web',
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
})
