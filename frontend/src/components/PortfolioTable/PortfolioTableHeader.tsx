import styles from './PortfolioTable.module.css';

interface PortfolioTableHeaderProps {
  isReadOnly: boolean;
}

export const PortfolioTableHeader = ({ isReadOnly }: PortfolioTableHeaderProps) => (
  <thead>
    <tr>
      <th className={styles.headerCell}>Asset</th>
      <th className={`${styles.headerCell} ${styles.headerCellRight}`}>Quantity</th>
      <th className={`${styles.headerCell} ${styles.headerCellRight}`}>Avg Price</th>
      <th className={`${styles.headerCell} ${styles.headerCellRight}`}>Current Price</th>
      <th className={`${styles.headerCell} ${styles.headerCellRight}`}>Market Value</th>
      <th className={`${styles.headerCell} ${styles.headerCellRight}`}>P&L</th>
      {!isReadOnly && (
        <th className={`${styles.headerCell} ${styles.headerCellCenter}`}>Actions</th>
      )}
    </tr>
  </thead>
);
