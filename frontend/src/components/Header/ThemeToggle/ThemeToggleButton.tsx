import { Button } from '@/components/shared/Button';
import React from 'react';
import styles from './ThemeToggle.module.css';
import clsx from 'clsx';

interface ThemeToggleButtonProps {
  active: boolean;
  onClick: () => void;
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}

export const ThemeToggleButton = ({
  active,
  onClick,
  icon: Icon,
  label,
  value,
}: ThemeToggleButtonProps) => {
  return (
    <Button
      data-testid={`theme-toggle-${value}`}
      variant="unstyled"
      size="unstyled"
      onClick={onClick}
      className={clsx(styles.toggleButton, active ? styles.active : styles.inactive)}
      aria-label={label}
      aria-pressed={active}
    >
      <Icon />
    </Button>
  );
};
