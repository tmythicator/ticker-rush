import * as React from 'react';
import styles from './AssetFininfoCell.module.css';

export interface AssetFininfoCellProps extends React.TdHTMLAttributes<HTMLTableCellElement> {
  variant?: 'default' | 'muted' | 'medium' | 'bold';
  trend?: 'up' | 'down' | 'neutral';
  ref?: React.Ref<HTMLTableCellElement>;
}

export const AssetFininfoCell = ({
  className,
  variant = 'default',
  trend,
  ref,
  ...props
}: AssetFininfoCellProps) => {
  return (
    <td
      ref={ref}
      className={`${styles.cell} ${className || ''}`}
      data-variant={variant}
      data-trend={trend}
      {...props}
    />
  );
};
