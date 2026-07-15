import { Button } from '@/components/shared/Button';
import React from 'react';
import styles from './ThemeToggle.module.css';

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
      className={styles.toggleButton}
      data-active={active}
      aria-label={label}
    >
      <Icon />
    </Button>
  );
};
