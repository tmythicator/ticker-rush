import type { PortfolioItem } from '@/types';
import { AssetInfoCell } from '@/components/PortfolioTable/AssetInfoCell';
import { RowActionsCell } from '@/components/PortfolioTable/RowActionsCell';
import { AssetFininfoCell } from '@/components/PortfolioTable/AssetFininfoCell';
import { usePortfolioRowState } from '@/hooks/usePortfolioRowState';

interface PortfolioRowProps {
  item: PortfolioItem;
  onSellClick: (item: PortfolioItem) => void;
  onTradeClick: (symbol: string) => void;
  isReadOnly?: boolean;
}

export const PortfolioRow = ({
  item,
  onSellClick,
  onTradeClick,
  isReadOnly = false,
}: PortfolioRowProps) => {
  const { symbol, source, isMarketClosed, quote, marketValue, pnl, pnlColorClass } =
    usePortfolioRowState(item);

  const loadingIndicator = <span className="animate-pulse">...</span>;

  return (
    <tr className="hover:bg-muted/50 transition-colors">
      <AssetInfoCell symbol={symbol} source={source} />

      <AssetFininfoCell variant="muted">{item.quantity}</AssetFininfoCell>

      <AssetFininfoCell variant="muted">${item.average_price.toFixed(2)}</AssetFininfoCell>

      <AssetFininfoCell variant="medium">
        {quote ? `$${quote.price.toFixed(2)}` : loadingIndicator}
      </AssetFininfoCell>

      <AssetFininfoCell variant="bold">{marketValue ?? loadingIndicator}</AssetFininfoCell>

      <AssetFininfoCell variant="bold" className={pnlColorClass}>
        {pnl ?? loadingIndicator}
      </AssetFininfoCell>

      {!isReadOnly && (
        <RowActionsCell
          item={item}
          isMarketClosed={isMarketClosed}
          onSellClick={onSellClick}
          onTradeClick={onTradeClick}
        />
      )}
    </tr>
  );
};
