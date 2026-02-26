import { promises as fs } from 'fs';
import path from 'path';

export async function generateSitemapAndRobots(outDir: string, siteUrl: string) {
  if (!siteUrl) {
    console.warn('VITE_SITE_URL is not defined! Skipping sitemap and robots generation.');
    return;
  }

  const baseUrl = siteUrl.endsWith('/') ? siteUrl : `${siteUrl}/`;
  console.log(`Generating SEO files for ${baseUrl}...`);

  await Promise.all([generateSitemap(outDir, baseUrl), generateRobots(outDir, baseUrl)]);

  console.log(`SEO files generated successfully in ${outDir}`);
}

async function generateSitemap(outDir: string, baseUrl: string) {
  const routes = [
    { url: '', changefreq: 'daily', priority: 1.0 },
    { url: 'login', changefreq: 'monthly', priority: 0.8 },
    { url: 'register', changefreq: 'monthly', priority: 0.8 },
    { url: 'leaderboard', changefreq: 'hourly', priority: 0.9 },
    { url: 'impressum', changefreq: 'yearly', priority: 0.3 },
    { url: 'privacy', changefreq: 'yearly', priority: 0.3 },
    { url: 'agb', changefreq: 'yearly', priority: 0.3 },
  ];

  const sitemapXml = `<?xml version="1.0" encoding="UTF-8"?>
  <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes
  .map(
    (route) => `  <url>
    <loc>${baseUrl}${route.url}</loc>
    <lastmod>${new Date().toISOString()}</lastmod>
    <changefreq>${route.changefreq}</changefreq>
    <priority>${route.priority.toFixed(1)}</priority>
    </url>`,
  )
  .join('\n')}
  </urlset>`;

  await fs.writeFile(path.join(outDir, 'sitemap.xml'), sitemapXml, 'utf-8');
}

async function generateRobots(outDir: string, baseUrl: string) {
  const robotsTxt = `User-agent: *
Allow: /

# Block search engines from indexing the private app areas
Disallow: /trade
Disallow: /profile

Sitemap: ${baseUrl}sitemap.xml
`;

  await fs.writeFile(path.join(outDir, 'robots.txt'), robotsTxt, 'utf-8');
}
