import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Header } from './Header';
import { useAuth } from '@/hooks/useAuth';
import { mockUserParticipating } from '@/test/mocks';
import styles from './Header.module.css';

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

describe('Header', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const renderHeader = () =>
    render(
      <MemoryRouter>
        <Header />
      </MemoryRouter>,
    );

  it('renders unauthenticated state buttons', () => {
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: false,
      user: null,
      logout: vi.fn(),
    });

    renderHeader();

    expect(screen.getByTestId('header-logo')).toBeInTheDocument();
    expect(screen.getByTestId('login-link')).toBeInTheDocument();
    expect(screen.getByTestId('register-link')).toBeInTheDocument();
  });

  it('renders authenticated layout with desktop and mobile toggles', () => {
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: true,
      user: mockUserParticipating,
      logout: vi.fn(),
    });

    renderHeader();

    const desktopLogout = screen.getByTestId('logout-button');
    const mobileMenuToggle = screen.getByTestId('mobile-menu-toggle');

    expect(desktopLogout.className).toContain(styles.desktopOnly);
    expect(mobileMenuToggle.className).toContain(styles.mdHide);
  });

  it('opens and closes mobile menu overlay on toggle click', () => {
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: true,
      user: mockUserParticipating,
      logout: vi.fn(),
    });

    renderHeader();

    // Mobile menu overlay should not be visible initially
    expect(screen.queryByTestId('mobile-menu')).not.toBeInTheDocument();

    const mobileMenuToggle = screen.getByTestId('mobile-menu-toggle');

    // Open menu
    fireEvent.click(mobileMenuToggle);
    expect(screen.getByTestId('mobile-menu')).toBeInTheDocument();

    // Close menu
    fireEvent.click(mobileMenuToggle);
    expect(screen.queryByTestId('mobile-menu')).not.toBeInTheDocument();
  });
});
