const extractColorFromCssVar = (varName: string, fallback: string, alpha?: number): string => {
  if (typeof window === 'undefined') return fallback;

  const root = document.documentElement;
  const value = getComputedStyle(root).getPropertyValue(varName).trim();

  if (!value) return fallback;
  // Handle HSL format (H S L)
  if (value.includes(' ')) {
    return alpha !== undefined ? `hsl(${value} / ${alpha})` : `hsl(${value})`;
  }

  return value;
};

export const getChartColors = () => {
  return {
    bgColor: extractColorFromCssVar('--background', '#ffffff'),
    textColor: extractColorFromCssVar('--foreground', '#000000'),
    borderColor: extractColorFromCssVar('--border', '#e2e8f0'),
    primaryColor: extractColorFromCssVar('--primary', '#10b981'),

    // Series colors derived from --primary with optional opacity
    areaTopColor: extractColorFromCssVar('--primary', '#26a69a8f', 0.56),
    areaBottomColor: extractColorFromCssVar('--primary', '#26a69a0a', 0.04),
    areaLineColor: extractColorFromCssVar('--primary', '#26a69a'),
  };
};
