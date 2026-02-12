import { useTheme } from 'next-themes';
import { useEffect } from 'react';

export const useThemeObserver = (callback: () => void) => {
  const { resolvedTheme } = useTheme();

  useEffect(() => {
    // When resolvedTheme changes, trigger the callback with a slight delay
    // to ensure the DOM has updated and styles are applied.
    const timer = setTimeout(() => {
      callback();
    }, 0);

    return () => clearTimeout(timer);
  }, [resolvedTheme, callback]);
};
