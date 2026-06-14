import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PublicProfilePage } from '@/pages/PublicProfilePage';
import { getPublicProfile } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockUserParticipating as mockUser } from '@/test/mocks';

vi.mock('@/lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/lib/api')>()),
  getPublicProfile: vi.fn(),
}));

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
      expect(screen.getByTestId('profile-name')).toHaveTextContent(
        mockUser.first_name + ' ' + mockUser.last_name,
      );
      expect(screen.getByTestId('profile-username')).toHaveTextContent('@' + mockUser.username);
      expect(screen.getByTestId('portfolio-row-aapl')).toBeInTheDocument();
    });
  });

  it('renders error for private/missing profile', async () => {
    (getPublicProfile as unknown as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { status: 404 },
    });
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('profile-unavailable')).toBeInTheDocument();
    });
  });
});
