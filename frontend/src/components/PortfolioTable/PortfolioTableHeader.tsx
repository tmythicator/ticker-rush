import styles from './PortfolioTable.module.css';

interface PortfolioTableHeaderProps {
  isReadOnly: boolean;
}

export const PortfolioTableHeader = ({ isReadOnly }: PortfolioTableHeaderProps) => (
  <thead>
    <tr>
      <th className={styles.headerCell}>Asset</th>
      <th className={styles.headerCell} data-align="right">
        Quantity
      </th>
      <th className={styles.headerCell} data-align="right">
        Avg Price
      </th>
      <th className={styles.headerCell} data-align="right">
        Current Price
      </th>
      <th className={styles.headerCell} data-align="right">
        Market Value
      </th>
      <th className={styles.headerCell} data-align="right">
        P&L
      </th>
      {!isReadOnly && (
        <th className={styles.headerCell} data-align="center">
          Actions
        </th>
      )}
    </tr>
  </thead>
);
