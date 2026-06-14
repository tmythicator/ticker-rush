import { Button } from '@/components/shared/Button';
import React from 'react';

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
      className={`relative z-10 flex h-full w-7 items-center justify-center transition-colors ${
        active
          ? 'bg-primary text-primary-foreground'
          : 'text-muted-foreground hover:text-foreground'
      }`}
      aria-label={label}
    >
      <Icon className="h-[14px] w-[14px]" />
    </Button>
  );
};
