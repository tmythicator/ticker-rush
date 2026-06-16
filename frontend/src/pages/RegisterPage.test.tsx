import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RegisterPage } from './RegisterPage';
import { register as apiRegister, login as apiLogin } from '@/lib/api';
import { useAuth } from '@/hooks/useAuth';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockUserParticipating } from '@/test/mocks';

vi.mock('@/lib/api', () => ({
  register: vi.fn(),
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

describe('RegisterPage', () => {
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
          <RegisterPage />
        </MemoryRouter>
      </QueryClientProvider>,
    );

  it('shows validation errors when fields are empty or invalid', async () => {
    const user = userEvent.setup();
    renderPage();

    const submitBtn = screen.getByTestId('register-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      const fieldErrors = screen.getAllByTestId('field-error');
      expect(fieldErrors).toHaveLength(5);
      expect(fieldErrors[0]).toHaveTextContent('Username must be at least 3 characters long');
      expect(fieldErrors[1]).toHaveTextContent('First name is required');
      expect(fieldErrors[2]).toHaveTextContent('Last name is required');
      expect(fieldErrors[3]).toHaveTextContent('Password must be at least 8 characters long');
      expect(fieldErrors[4]).toHaveTextContent(
        'You must accept the Terms and Conditions and Privacy Policy',
      );
    });
  });

  it('renders backend error using ErrorMessage when registration fails', async () => {
    const user = userEvent.setup();
    (apiRegister as unknown as ReturnType<typeof vi.fn>).mockRejectedValue(
      new Error('Username already taken'),
    );

    renderPage();

    await user.type(screen.getByTestId('username-input'), 'existinguser');
    await user.type(screen.getByTestId('first-name-input'), 'John');
    await user.type(screen.getByTestId('last-name-input'), 'Doe');
    await user.type(screen.getByTestId('password-input'), 'password123');
    await user.click(screen.getByTestId('agb-checkbox'));

    const submitBtn = screen.getByTestId('register-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      const errorAlert = screen.getByTestId('error-message');
      expect(errorAlert).toBeInTheDocument();
      expect(errorAlert).toHaveTextContent('Username already taken');
    });
  });

  it('redirects and updates auth state on successful registration', async () => {
    const user = userEvent.setup();
    (apiRegister as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);
    (apiLogin as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);

    renderPage();

    await user.type(screen.getByTestId('username-input'), 'newuser');
    await user.type(screen.getByTestId('first-name-input'), 'John');
    await user.type(screen.getByTestId('last-name-input'), 'Doe');
    await user.type(screen.getByTestId('password-input'), 'password123');
    await user.click(screen.getByTestId('agb-checkbox'));

    const submitBtn = screen.getByTestId('register-submit');
    await user.click(submitBtn);

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith(mockUserParticipating);
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});
