import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ErrorMessage } from './ErrorMessage';

describe('ErrorMessage', () => {
  it('renders message string correctly', () => {
    render(<ErrorMessage message="Something went wrong" />);
    const alert = screen.getByRole('alert');
    expect(alert).toBeInTheDocument();
    expect(alert.textContent).toBe('Something went wrong');
  });

  it('renders children correctly', () => {
    render(<ErrorMessage>Children error message</ErrorMessage>);
    const alert = screen.getByRole('alert');
    expect(alert).toBeInTheDocument();
    expect(alert.textContent).toBe('Children error message');
  });

  it('does not render anything if no message or children are provided', () => {
    const { container } = render(<ErrorMessage />);
    expect(container.firstChild).toBeNull();
  });

  it('applies variant classes correctly', () => {
    const { rerender } = render(<ErrorMessage message="Error" variant="sm" />);
    let alert = screen.getByRole('alert');
    expect(alert.className).toContain('p-3');
    expect(alert.className).toContain('text-sm');

    rerender(<ErrorMessage message="Error" variant="xs" />);
    alert = screen.getByRole('alert');
    expect(alert.className).toContain('p-2');
    expect(alert.className).toContain('text-xs');
  });

  it('applies custom className and forwards native div properties', () => {
    render(<ErrorMessage message="Error" className="custom-class" data-testid="custom-error" />);
    const alert = screen.getByTestId('custom-error');
    expect(alert.className).toContain('custom-class');
    expect(alert).toHaveAttribute('role', 'alert');
  });
});
