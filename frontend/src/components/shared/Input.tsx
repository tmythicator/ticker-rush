import * as React from 'react';
import { type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { inputVariants } from './inputVariants';

export interface InputProps
  extends Omit<React.ComponentProps<'input'>, 'size'>, VariantProps<typeof inputVariants> {
  ref?: React.Ref<HTMLInputElement>;
}

export const Input = ({ className, type, variant, size, ref, ...props }: InputProps) => {
  return (
    <input
      type={type}
      className={cn(inputVariants({ variant, size, className }))}
      ref={ref}
      {...props}
    />
  );
};
