import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ThemeToggle } from './ThemeToggle';
import { useIsMounted } from '@/hooks/useIsMounted';
import { useTheme } from 'next-themes';

vi.mock('@/hooks/useIsMounted', () => ({
  useIsMounted: vi.fn(),
}));

vi.mock('next-themes', () => ({
  useTheme: vi.fn(),
}));

describe('ThemeToggle', () => {
  const mockSetTheme = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useIsMounted as unknown as ReturnType<typeof vi.fn>).mockReturnValue(true);
    (useTheme as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      theme: 'system',
      setTheme: mockSetTheme,
    });
  });

  it('renders loading state when not mounted', () => {
    (useIsMounted as unknown as ReturnType<typeof vi.fn>).mockReturnValue(false);
    render(<ThemeToggle />);

    expect(screen.queryByTestId('theme-toggle-dark')).not.toBeInTheDocument();
  });

  it('renders theme buttons when mounted', () => {
    render(<ThemeToggle />);

    expect(screen.getByTestId('theme-toggle-light')).toBeInTheDocument();
    expect(screen.getByTestId('theme-toggle-system')).toBeInTheDocument();
    expect(screen.getByTestId('theme-toggle-dark')).toBeInTheDocument();
  });

  it('highlights the active theme button', () => {
    (useTheme as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      theme: 'dark',
      setTheme: mockSetTheme,
    });

    render(<ThemeToggle />);

    const darkButton = screen.getByTestId('theme-toggle-dark');
    const lightButton = screen.getByTestId('theme-toggle-light');

    expect(darkButton.className).toContain('bg-primary');
    expect(lightButton.className).not.toContain('bg-primary');
    expect(lightButton.className).toContain('text-muted-foreground');
  });

  it('calls setTheme with selected theme on click', async () => {
    const user = userEvent.setup();
    render(<ThemeToggle />);

    const darkButton = screen.getByTestId('theme-toggle-dark');
    await user.click(darkButton);

    expect(mockSetTheme).toHaveBeenCalledWith('dark');
  });
});
