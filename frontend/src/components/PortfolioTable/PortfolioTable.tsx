import type { PortfolioItem } from '@/types';
import { PortfolioRow } from './PortfolioRow';
import { PortfolioTableHeader } from './PortfolioTableHeader';
import { TableEmptyState } from './TableEmptyState';

interface PortfolioTableProps {
  items: PortfolioItem[];
  isReadOnly?: boolean;
  onSellClick?: (item: PortfolioItem) => void;
  onTradeClick?: (symbol: string) => void;
}

export const PortfolioTable = ({
  items,
  isReadOnly = false,
  onSellClick,
  onTradeClick,
}: PortfolioTableProps) => {
  const isEmpty = items.length === 0;

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left" data-testid="portfolio-table">
        <PortfolioTableHeader isReadOnly={isReadOnly} />
        <tbody className="divide-y divide-border">
          {items.map((item) => (
            <PortfolioRow
              key={item.stock_symbol}
              item={item}
              onSellClick={onSellClick || (() => {})}
              onTradeClick={onTradeClick || (() => {})}
              isReadOnly={isReadOnly}
            />
          ))}
          {isEmpty && <TableEmptyState isReadOnly={isReadOnly} />}
        </tbody>
      </table>
    </div>
  );
};
