interface TableEmptyStateProps {
  isReadOnly: boolean;
}

export const TableEmptyState = ({ isReadOnly }: TableEmptyStateProps) => (
  <tr>
    <td
      colSpan={isReadOnly ? 6 : 7}
      data-testid="portfolio-empty-state"
      className="px-6 py-12 text-center italic text-muted-foreground"
    >
      No assets found in your portfolio.{!isReadOnly && ' Start trading!'}
    </td>
  </tr>
);
