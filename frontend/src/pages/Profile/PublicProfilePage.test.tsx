import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PublicProfilePage } from './PublicProfilePage';
import { getPublicProfile } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockPublicProfile } from '@/test/mocks';

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
    vi.mocked(getPublicProfile).mockResolvedValue(mockPublicProfile);
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('profile-name')).toHaveTextContent(
        `${mockPublicProfile.first_name} ${mockPublicProfile.last_name}`.trim(),
      );

      expect(screen.getByTestId('profile-username')).toHaveTextContent(
        `@${mockPublicProfile.username}`,
      );

      if (mockPublicProfile.website) {
        const link = screen.getByRole('link', {
          name: new RegExp(`Visit website: ${mockPublicProfile.website}`, 'i'),
        });
        expect(link).toHaveAttribute('href', mockPublicProfile.website);
        expect(link).toHaveAttribute('target', '_blank');
      }

      expect(screen.getByTestId('portfolio-row-aapl')).toBeInTheDocument();
    });
  });

  it('renders error for private/missing profile', async () => {
    vi.mocked(getPublicProfile).mockRejectedValue({
      response: { status: 404 },
    });
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('profile-unavailable')).toBeInTheDocument();
    });
  });
});
