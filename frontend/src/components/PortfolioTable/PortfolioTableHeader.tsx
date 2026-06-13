interface PortfolioTableHeaderProps {
  isReadOnly: boolean;
}

export const PortfolioTableHeader = ({ isReadOnly }: PortfolioTableHeaderProps) => (
  <thead>
    <tr className="bg-muted text-muted-foreground text-xs uppercase tracking-wider font-bold">
      <th className="px-6 py-4">Asset</th>
      <th className="px-6 py-4 text-right">Quantity</th>
      <th className="px-6 py-4 text-right">Avg Price</th>
      <th className="px-6 py-4 text-right">Current Price</th>
      <th className="px-6 py-4 text-right">Market Value</th>
      <th className="px-6 py-4 text-right">P&L</th>
      {!isReadOnly && <th className="px-6 py-4 text-right">Actions</th>}
    </tr>
  </thead>
);
