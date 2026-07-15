import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ErrorMessage } from './ErrorMessage';
import styles from './ErrorMessage.module.css';

describe('ErrorMessage', () => {
  it('renders message string correctly', () => {
    render(<ErrorMessage data-testid="error" message="Something went wrong" />);
    const alert = screen.getByTestId('error');
    expect(alert).toBeInTheDocument();
    expect(alert.textContent).toBe('Something went wrong');
  });

  it('renders children correctly', () => {
    render(<ErrorMessage data-testid="error">Children error message</ErrorMessage>);
    const alert = screen.getByTestId('error');
    expect(alert).toBeInTheDocument();
    expect(alert.textContent).toBe('Children error message');
  });

  it('does not render anything if no message or children are provided', () => {
    const { container } = render(<ErrorMessage />);
    expect(container.firstChild).toBeNull();
  });

  it('applies variant classes correctly', () => {
    const { rerender } = render(<ErrorMessage message="Error" variant="sm" data-testid="error" />);
    let alert = screen.getByTestId('error');
    expect(alert.className).toContain(styles.sizeSm);

    rerender(<ErrorMessage message="Error" variant="xs" data-testid="error" />);
    alert = screen.getByTestId('error');
    expect(alert.className).toContain(styles.sizeXs);
  });

  it('applies custom className and forwards native div properties', () => {
    render(<ErrorMessage message="Error" className="custom-class" data-testid="custom-error" />);
    const alert = screen.getByTestId('custom-error');
    expect(alert.className).toContain('custom-class');
    expect(alert).toHaveAttribute('role', 'alert');
  });
});
