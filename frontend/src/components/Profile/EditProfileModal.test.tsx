import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { EditProfileModal } from './EditProfileModal';
import { updateUser as apiUpdateUser, getUser as apiGetUser, deleteUser as apiDeleteUser } from '@/lib/api';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockUserParticipating } from '@/test/mocks';

import { useAuth } from '@/hooks/useAuth';

vi.mock('@/lib/api', () => ({
  updateUser: vi.fn(),
  getUser: vi.fn(),
  deleteUser: vi.fn(),
}));

vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

describe('EditProfileModal', () => {
  let queryClient: QueryClient;
  const mockOnClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    (apiGetUser as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mockUserParticipating);
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: mockUserParticipating,
      setUser: vi.fn(),
      login: vi.fn(),
    });
  });

  const renderModal = (isOpen = true) =>
    render(
      <QueryClientProvider client={queryClient}>
        <EditProfileModal isOpen={isOpen} onClose={mockOnClose} />
      </QueryClientProvider>,
    );

  it('renders modal contents when open', async () => {
    renderModal();

    expect(screen.getByRole('heading', { name: /edit profile/i })).toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByTestId('first-name-input')).toHaveValue('Test');
      expect(screen.getByTestId('last-name-input')).toHaveValue('User');
      expect(screen.getByTestId('website-input')).toHaveValue('https://example.com');
    });
  });

  it('does not render modal contents when closed', () => {
    renderModal(false);
    expect(screen.queryByRole('heading', { name: /edit profile/i })).not.toBeInTheDocument();
  });

  it('renders backend error using ErrorMessage when profile update fails', async () => {
    const user = userEvent.setup();
    (apiUpdateUser as unknown as ReturnType<typeof vi.fn>).mockRejectedValue(
      new Error('Failed to save profile changes'),
    );

    renderModal();

    await waitFor(() => {
      expect(screen.getByTestId('first-name-input')).toHaveValue('Test');
    });

    await user.clear(screen.getByTestId('first-name-input'));
    await user.type(screen.getByTestId('first-name-input'), 'NewName');

    const saveBtn = screen.getByTestId('edit-profile-submit');
    await user.click(saveBtn);

    await waitFor(() => {
      const errorAlert = screen.getByTestId('error-message');
      expect(errorAlert).toBeInTheDocument();
      expect(errorAlert).toHaveTextContent('Failed to save profile changes');
    });
  });

  it('calls onClose and updates profile on successful submit', async () => {
    const user = userEvent.setup();
    (apiUpdateUser as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
      ...mockUserParticipating,
      first_name: 'NewName',
    });

    renderModal();

    await waitFor(() => {
      expect(screen.getByTestId('first-name-input')).toHaveValue('Test');
    });

    await user.clear(screen.getByTestId('first-name-input'));
    await user.type(screen.getByTestId('first-name-input'), 'NewName');

    const saveBtn = screen.getByTestId('edit-profile-submit');
    await user.click(saveBtn);

    await waitFor(() => {
      expect(apiUpdateUser).toHaveBeenCalledWith({
        first_name: 'NewName',
        last_name: 'User',
        website: 'https://example.com',
        is_public: true,
        update_mask: { paths: ['first_name'] },
      });
      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  it('displays validation error for invalid website URLs', async () => {
    const user = userEvent.setup();
    renderModal();

    await waitFor(() => {
      expect(screen.getByTestId('website-input')).toHaveValue('https://example.com');
    });

    // 1. Check data URL scheme restriction
    await user.clear(screen.getByTestId('website-input'));
    await user.type(screen.getByTestId('website-input'), 'data:text/html,hello');
    await user.click(screen.getByTestId('edit-profile-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('field-error')).toHaveTextContent(
        'Website must start with http:// or https://',
      );
    });

    // 2. Check missing protocol restriction
    await user.clear(screen.getByTestId('website-input'));
    await user.type(screen.getByTestId('website-input'), 'google.com');
    await user.click(screen.getByTestId('edit-profile-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('field-error')).toHaveTextContent('Invalid URL format');
    });

    // 3. Check invalid URL format with correct protocol
    await user.clear(screen.getByTestId('website-input'));
    await user.type(screen.getByTestId('website-input'), 'http://invalid url');
    await user.click(screen.getByTestId('edit-profile-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('field-error')).toHaveTextContent('Invalid URL format');
    });

    // 4. Check website length restriction (>200 chars)
    await user.clear(screen.getByTestId('website-input'));
    const longWebsite = 'https://' + 'a'.repeat(200);
    await user.type(screen.getByTestId('website-input'), longWebsite);
    await user.click(screen.getByTestId('edit-profile-submit'));

    await waitFor(() => {
      expect(screen.getByTestId('field-error')).toHaveTextContent(
        'Website must be at most 200 characters long',
      );
    });
  });

  it('calls deleteUser and logs out user on deletion confirmation', async () => {
    const user = userEvent.setup();
    const mockLogin = vi.fn();
    (useAuth as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: mockUserParticipating,
      setUser: vi.fn(),
      login: mockLogin,
    });
    (apiDeleteUser as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    renderModal();

    const deleteBtn = screen.getByTestId('delete-profile-button');
    await user.click(deleteBtn);

    // Verify custom modal is opened via its confirm button
    const confirmBtn = screen.getByTestId('delete-profile-confirm-submit');
    expect(confirmBtn).toBeInTheDocument();
    await user.click(confirmBtn);

    await waitFor(() => {
      expect(apiDeleteUser).toHaveBeenCalled();
      expect(mockLogin).toHaveBeenCalledWith(null);
      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  it('does not call deleteUser if deletion is not confirmed', async () => {
    const user = userEvent.setup();
    (apiDeleteUser as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    renderModal();

    const deleteBtn = screen.getByTestId('delete-profile-button');
    await user.click(deleteBtn);

    // Verify custom modal is opened via its cancel button
    const cancelBtn = screen.getByTestId('delete-profile-confirm-cancel');
    expect(cancelBtn).toBeInTheDocument();
    await user.click(cancelBtn);

    // Verify custom modal is closed
    await waitFor(() => {
      expect(screen.queryByTestId('delete-profile-confirm-cancel')).not.toBeInTheDocument();
    });

    expect(apiDeleteUser).not.toHaveBeenCalled();
  });

  it('renders delete error message when deletion fails', async () => {
    const user = userEvent.setup();
    (apiDeleteUser as unknown as ReturnType<typeof vi.fn>).mockRejectedValue(
      new Error('Failed to delete profile'),
    );

    renderModal();

    const deleteBtn = screen.getByTestId('delete-profile-button');
    await user.click(deleteBtn);

    const confirmBtn = screen.getByTestId('delete-profile-confirm-submit');
    await user.click(confirmBtn);

    await waitFor(() => {
      const errorAlert = screen.getByTestId('error-message');
      expect(errorAlert).toBeInTheDocument();
      expect(errorAlert).toHaveTextContent('Failed to delete profile');
    });
  });
});
