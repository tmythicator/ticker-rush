import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import App from './App';
import { useUserQuery } from '@/hooks/useUserQuery';
import { mockUserParticipating } from '@/test/mocks';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

vi.mock('@/components/Home/HomeChart', () => ({
  HomeChart: () => <div data-testid="mock-home-chart" />,
}));

vi.mock('@/hooks/useUserQuery', () => ({
  useUserQuery: vi.fn(),
}));

describe('App', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });
  });

  const renderApp = () =>
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>,
    );

  it('renders landing page with login and register links when unauthenticated', async () => {
    (useUserQuery as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: null,
      isLoading: false,
    });

    renderApp();

    expect(screen.getByTestId('app-header')).toBeInTheDocument();
    expect(screen.getByTestId('app-footer')).toBeInTheDocument();
    expect(screen.getByTestId('login-link')).toBeInTheDocument();
    expect(screen.getByTestId('register-link')).toBeInTheDocument();
    expect(screen.queryByTestId('logout-button')).not.toBeInTheDocument();
  });

  it('renders dashboard shell and logout button when authenticated', async () => {
    (useUserQuery as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockUserParticipating,
      isLoading: false,
    });

    renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('app-header')).toBeInTheDocument();
      expect(screen.getByTestId('logout-button')).toBeInTheDocument();
      expect(screen.queryByTestId('login-link')).not.toBeInTheDocument();
      expect(screen.queryByTestId('register-link')).not.toBeInTheDocument();
    });
  });
});
