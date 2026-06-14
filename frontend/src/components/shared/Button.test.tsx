import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { Button } from './Button';

describe('Button', () => {
  it('renders children correctly', () => {
    render(<Button data-testid="test-button">Click Me</Button>);
    const button = screen.getByTestId('test-button');
    expect(button).toBeInTheDocument();
    expect(button.textContent).toBe('Click Me');
  });

  it('applies variant classes correctly', () => {
    const { rerender } = render(
      <Button data-testid="test-button" variant="outline">
        Outline
      </Button>,
    );
    let button = screen.getByTestId('test-button');
    expect(button.className).toContain('bg-background');

    rerender(
      <Button data-testid="test-button" variant="destructive">
        Destructive
      </Button>,
    );
    button = screen.getByTestId('test-button');
    expect(button.className).toContain('bg-destructive');
  });

  it('applies custom className and native attributes', () => {
    render(
      <Button data-testid="test-button" className="custom-class" type="submit" disabled>
        Submit
      </Button>,
    );
    const button = screen.getByTestId('test-button');
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute('type', 'submit');
    expect(button.className).toContain('custom-class');
  });

  it('triggers onClick handler when clicked', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(
      <Button data-testid="test-button" onClick={handleClick}>
        Clickable
      </Button>,
    );

    const button = screen.getByTestId('test-button');
    await user.click(button);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('does not trigger onClick when disabled', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(
      <Button data-testid="test-button" onClick={handleClick} disabled>
        Disabled
      </Button>,
    );

    const button = screen.getByTestId('test-button');
    await user.click(button);

    expect(handleClick).not.toHaveBeenCalled();
  });

  it('forwards refs correctly', () => {
    const ref = React.createRef<HTMLButtonElement>();
    render(
      <Button data-testid="test-button" ref={ref}>
        With Ref
      </Button>,
    );

    expect(ref.current).toBeInstanceOf(HTMLButtonElement);
    expect(ref.current?.getAttribute('data-testid')).toBe('test-button');
  });
});
