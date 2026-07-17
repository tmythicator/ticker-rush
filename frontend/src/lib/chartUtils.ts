const extractColorFromCssVar = (varName: string, fallback: string, alpha?: number): string => {
  if (typeof window === 'undefined') return fallback;

  const root = document.documentElement;
  const value = getComputedStyle(root).getPropertyValue(varName).trim();

  if (!value) return fallback;

  // Check if it already has hsl/hsla wrapper
  const match = value.match(/^(hsl|hsla|rgb|rgba)\(([^)]+)\)$/i);
  if (match) {
    if (alpha === undefined) return value;

    // 'hsla' -> 'hsl', 'rgba' -> 'rgb'
    const basePrefix = match[1].toLowerCase().replace('a', '');
    const inner = match[2].trim();

    // Format
    return inner.includes(',')
      ? `${basePrefix}a(${inner}, ${alpha})`
      : `${basePrefix}(${inner} / ${alpha})`;
  }

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

/**
 * Formats a timestamp to a time string.
 * @param timestamp The timestamp to format.
 * @returns The formatted time string.
 */
export const formatTime = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
