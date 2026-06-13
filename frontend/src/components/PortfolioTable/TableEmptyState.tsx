interface TableEmptyStateProps {
  isReadOnly: boolean;
}

export const TableEmptyState = ({ isReadOnly }: TableEmptyStateProps) => (
  <tr>
    <td
      colSpan={isReadOnly ? 6 : 7}
      className="px-6 py-12 text-center text-muted-foreground italic"
    >
      No assets found in your portfolio.{!isReadOnly && ' Start trading!'}
    </td>
  </tr>
);
