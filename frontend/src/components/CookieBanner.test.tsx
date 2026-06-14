import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, beforeEach } from 'vitest';
import { CookieBanner } from './CookieBanner';

describe('CookieBanner', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('renders cookie banner when no consent is in localStorage', () => {
    render(<CookieBanner />);

    expect(screen.getByTestId('cookie-banner')).toBeInTheDocument();
    expect(screen.getByTestId('cookie-banner-accept-button')).toBeInTheDocument();
  });

  it('does not render when consent has already been given', () => {
    localStorage.setItem('cookie-consent', 'true');
    render(<CookieBanner />);

    expect(screen.queryByTestId('cookie-banner')).not.toBeInTheDocument();
  });

  it('sets consent in localStorage and hides banner when accepted', async () => {
    const user = userEvent.setup();
    render(<CookieBanner />);

    const button = screen.getByTestId('cookie-banner-accept-button');
    await user.click(button);

    expect(localStorage.getItem('cookie-consent')).toBe('true');
    expect(screen.queryByTestId('cookie-banner')).not.toBeInTheDocument();
  });
});
