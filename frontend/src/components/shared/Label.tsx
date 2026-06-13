import * as React from 'react';
import { type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { labelVariants } from './labelVariants';

export interface LabelProps
  extends React.LabelHTMLAttributes<HTMLLabelElement>, VariantProps<typeof labelVariants> {
  ref?: React.Ref<HTMLLabelElement>;
}

export const Label = ({ className, variant, ref, ...props }: LabelProps) => (
  <label ref={ref} className={cn(labelVariants({ variant, className }))} {...props} />
);
