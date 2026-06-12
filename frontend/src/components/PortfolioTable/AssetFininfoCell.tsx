import * as React from 'react';
import { type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { assetFininfoCellVariants } from './assetFininfoCellVariants';

export interface AssetFininfoCellProps
  extends
    React.TdHTMLAttributes<HTMLTableCellElement>,
    VariantProps<typeof assetFininfoCellVariants> {
  ref?: React.Ref<HTMLTableCellElement>;
}

export const AssetFininfoCell = ({ className, variant, ref, ...props }: AssetFininfoCellProps) => (
  <td ref={ref} className={cn(assetFininfoCellVariants({ variant, className }))} {...props} />
);
