import type { PortfolioItem } from '@/types';
import { AssetInfoCell } from '@/components/PortfolioTable/AssetInfoCell';
import { RowActionsCell } from '@/components/PortfolioTable/RowActionsCell';
import { AssetFininfoCell } from '@/components/PortfolioTable/AssetFininfoCell';
import { usePortfolioRowState } from '@/hooks/usePortfolioRowState';
import styles from './PortfolioTable.module.css';

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
  const { symbol, source, isMarketClosed, isTradable, quote, marketValue, pnl, pnlStatus } =
    usePortfolioRowState(item);

  const loadingIndicator = <span>...</span>;

  return (
    <tr data-testid={`portfolio-row-${symbol.toLowerCase()}`} className={styles.row}>
      <AssetInfoCell symbol={symbol} source={source} isTradable={isTradable} />

      <AssetFininfoCell variant="muted">{item.quantity}</AssetFininfoCell>

      <AssetFininfoCell variant="muted">${item.average_price.toFixed(2)}</AssetFininfoCell>

      <AssetFininfoCell variant="medium">
        {quote ? `$${quote.price.toFixed(2)}` : loadingIndicator}
      </AssetFininfoCell>

      <AssetFininfoCell variant="bold">{marketValue ?? loadingIndicator}</AssetFininfoCell>

      <AssetFininfoCell
        variant="bold"
        trend={pnlStatus === 'positive' ? 'up' : pnlStatus === 'negative' ? 'down' : 'neutral'}
      >
        {pnl ?? loadingIndicator}
      </AssetFininfoCell>

      {!isReadOnly && (
        <RowActionsCell
          item={item}
          isMarketClosed={isMarketClosed}
          isTradable={isTradable}
          onSellClick={onSellClick}
          onTradeClick={onTradeClick}
        />
      )}
    </tr>
  );
};
