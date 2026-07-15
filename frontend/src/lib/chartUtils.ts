const extractColorFromCssVar = (varName: string, fallback: string, alpha?: number): string => {
  if (typeof window === 'undefined') return fallback;

  const root = document.documentElement;
  const value = getComputedStyle(root).getPropertyValue(varName).trim();

  if (!value) return fallback;

  // Check if it already has hsl/hsla wrapper
  const hslMatch = value.match(/^(hsl|hsla)\(([^)]+)\)$/i);
  if (hslMatch) {
    const inner = hslMatch[2].trim();
    if (alpha !== undefined) {
      if (inner.includes(',')) {
        return `hsla(${inner}, ${alpha})`;
      } else {
        return `hsl(${inner} / ${alpha})`;
      }
    }
    return value;
  }

  // Check if it already has rgb/rgba wrapper
  const rgbMatch = value.match(/^(rgb|rgba)\(([^)]+)\)$/i);
  if (rgbMatch) {
    const inner = rgbMatch[2].trim();
    if (alpha !== undefined) {
      if (inner.includes(',')) {
        return `rgba(${inner}, ${alpha})`;
      } else {
        return `rgb(${inner} / ${alpha})`;
      }
    }
    return value;
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
