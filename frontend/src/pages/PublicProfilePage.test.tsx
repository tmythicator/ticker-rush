import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PublicProfilePage } from '@/pages/PublicProfilePage';
import { getPublicProfile } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the API
vi.mock('@/lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/lib/api')>()),
  getPublicProfile: vi.fn(),
}));

const mockUser = {
  id: 1,
  username: 'testuser',
  first_name: 'Test',
  last_name: 'User',
  website: 'https://example.com',
  is_public: true,
  balance: 5000.0,
  portfolio: { AAPL: { stock_symbol: 'AAPL', quantity: 10, average_price: 150.0 } },
};

const renderPage = () =>
  render(
    <QueryClientProvider
      client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}
    >
      <MemoryRouter initialEntries={['/users/testuser']}>
        <Routes>
          <Route path="/users/:username" element={<PublicProfilePage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );

describe('PublicProfilePage', () => {
  beforeEach(() => vi.clearAllMocks());

  it('renders public profile correctly', async () => {
    (getPublicProfile as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUser);
    renderPage();

    await waitFor(() => {
      expect(screen.getByText(mockUser.first_name + ' ' + mockUser.last_name)).toBeInTheDocument();
      expect(screen.getByText(`@${mockUser.username}`)).toBeInTheDocument();
      expect(screen.getByText(mockUser.portfolio.AAPL.stock_symbol)).toBeInTheDocument();
    });
  });

  it('renders error for private/missing profile', async () => {
    (getPublicProfile as unknown as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { status: 404 },
    });
    renderPage();

    await waitFor(() => {
      expect(screen.getByText('Profile Unavailable')).toBeInTheDocument();
    });
  });
});
