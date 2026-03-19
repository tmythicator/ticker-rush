import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

/// <reference types="vitest" />
import { promises as fs } from 'fs';
import { generateSitemapAndRobots } from './config/seo';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const rootEnv = loadEnv(mode, path.resolve(__dirname, '../'), '');
  const appEnv = loadEnv(mode, process.cwd(), '');
  const exchangePort = parseInt(rootEnv.SERVER_PORT) || 8081;
  const clientPort = parseInt(rootEnv.CLIENT_PORT) || 5173;

  return {
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
        '@icons': path.resolve(__dirname, './src/components/icons'),
      },
    },
    esbuild: {
      drop: mode === 'production' ? ['console', 'debugger'] : [],
    },
    define: {
      'import.meta.env.VITE_API_URL': JSON.stringify('/api'),
    },
    server: {
      port: clientPort,
      host: true,
      // Forward requests to backend (avoid CORS hell)
      proxy: {
        '/api': {
          target: `http://localhost:${exchangePort}`,
          changeOrigin: true,
          secure: false,
        },
      },
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './src/test/setup.ts',
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: [
              'react',
              'react-dom',
              'react-router-dom',
              '@tanstack/react-query',
              'class-variance-authority',
              'clsx',
              'tailwind-merge',
              'next-themes',
              '@radix-ui/react-label',
              '@radix-ui/react-slot',
            ],
            forms: ['react-hook-form', 'zod', '@hookform/resolvers'],
            charts: ['lightweight-charts'],
          },
        },
      },
    },
    plugins: [
      react(),
      {
        name: 'generate-seo-files',
        apply: 'build',
        buildStart: async () => {
          const publicDir = path.resolve(__dirname, 'public');
          await fs.mkdir(publicDir, { recursive: true });
          const siteUrl = appEnv.VITE_SITE_URL || 'https://example.com';
          await generateSitemapAndRobots(publicDir, siteUrl);
        },
      },
    ],
  };
});
