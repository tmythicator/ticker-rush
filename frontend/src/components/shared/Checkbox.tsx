import * as React from 'react';
import { cn } from '@/lib/utils';

export interface CheckboxProps extends React.ComponentProps<'input'> {
  ref?: React.Ref<HTMLInputElement>;
}

export const Checkbox = ({ className, ref, ...props }: CheckboxProps) => {
  return (
    <input
      type="checkbox"
      className={cn(
        'h-5 w-5 cursor-pointer rounded border-border bg-background text-primary accent-primary transition-colors focus:ring-primary/50',
        className,
      )}
      ref={ref}
      {...props}
    />
  );
};
