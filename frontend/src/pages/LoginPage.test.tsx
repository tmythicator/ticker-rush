import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { LoginPage } from './LoginPage';
import { login as apiLogin } from '@/lib/api';
import { useAuth } from '@/hooks/useAuth';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockUserParticipating } from '@/test/mocks';

vi.mock('@/lib/api', () => ({
  login: vi.fn(),
}));

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => ({
  ...(await importOriginal<typeof import('react-router-dom')>()),
  useNavigate: () => mockNavigate,
}));

describe('LoginPage', () => {
  let queryClient: QueryClient;
  const mockLogin = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      login: mockLogin,
    });
  });

  const renderPage = () =>
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <LoginPage />
        </MemoryRouter>
      </QueryClientProvider>,
    );

  it('shows validation errors when fields are empty', async () => {
    const user = userEvent.setup();
    renderPage();

    const submitBtn = screen.getByTestId('login-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      const fieldErrors = screen.getAllByTestId('field-error');
      expect(fieldErrors).toHaveLength(2);
      expect(fieldErrors[0]).toHaveTextContent('Username must be at least 3 characters long');
      expect(fieldErrors[1]).toHaveTextContent('Password must be at least 8 characters long');
    });
  });

  it('renders backend error using ErrorMessage when login fails', async () => {
    const user = userEvent.setup();
    (apiLogin as unknown as ReturnType<typeof vi.fn>).mockRejectedValue(
      new Error('Invalid username or password'),
    );

    renderPage();

    await user.type(screen.getByTestId('username-input'), 'wronguser');
    await user.type(screen.getByTestId('password-input'), 'wrongpassword');

    const submitBtn = screen.getByTestId('login-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      const errorAlert = screen.getByTestId('error-message');
      expect(errorAlert).toBeInTheDocument();
      expect(errorAlert).toHaveTextContent('Invalid username or password');
    });
  });

  it('redirects and updates auth state on successful login', async () => {
    const user = userEvent.setup();
    (apiLogin as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);

    renderPage();

    await user.type(screen.getByTestId('username-input'), 'testuser');
    await user.type(screen.getByTestId('password-input'), 'password123');

    const submitBtn = screen.getByTestId('login-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith(mockUserParticipating);
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});
