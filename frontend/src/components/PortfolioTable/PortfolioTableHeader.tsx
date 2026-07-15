import styles from './PortfolioTable.module.css';
import { cva } from 'class-variance-authority';

interface PortfolioTableHeaderProps {
  isReadOnly: boolean;
}

const headerVariants = cva(styles.headerCell, {
  variants: {
    align: {
      left: styles.alignLeft,
      center: styles.alignCenter,
      right: styles.alignRight,
    },
  },
  defaultVariants: {
    align: 'center',
  },
});

export const PortfolioTableHeader = ({ isReadOnly }: PortfolioTableHeaderProps) => (
  <thead>
    <tr>
      <th className={headerVariants()}>Asset</th>
      <th className={headerVariants()}>Quantity</th>
      <th className={headerVariants()}>Avg Price</th>
      <th className={headerVariants()}>Current Price</th>
      <th className={headerVariants()}>Market Value</th>
      <th className={headerVariants()}>P&L</th>
      {!isReadOnly && <th className={headerVariants()}>Actions</th>}
    </tr>
  </thead>
);
