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
        'h-5 w-5 rounded border-border text-primary focus:ring-primary/50 bg-background accent-primary cursor-pointer transition-colors',
        className,
      )}
      ref={ref}
      {...props}
    />
  );
};
