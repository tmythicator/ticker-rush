import styles from './PortfolioTable.module.css';

interface TableEmptyStateProps {
  isReadOnly: boolean;
}

export const TableEmptyState = ({ isReadOnly }: TableEmptyStateProps) => (
  <tr>
    <td
      colSpan={isReadOnly ? 6 : 7}
      data-testid="portfolio-empty-state"
      className={styles.emptyCell}
    >
      No assets found in your portfolio.{!isReadOnly && ' Start trading!'}
    </td>
  </tr>
);
