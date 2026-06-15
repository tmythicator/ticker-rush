import { IconArrowRight, IconTrash } from '@/components/icons/CustomIcons';
import { Button } from '@/components/shared/Button';
import type { PortfolioItem } from '@/types';

interface RowActionsCellProps {
  item: PortfolioItem;
  isMarketClosed: boolean;
  isTradable?: boolean;
  onSellClick: (item: PortfolioItem) => void;
  onTradeClick: (symbol: string) => void;
}

export const RowActionsCell = ({
  item,
  isMarketClosed,
  isTradable = true,
  onSellClick,
  onTradeClick,
}: RowActionsCellProps) => {
  const isActionDisabled = isMarketClosed || !isTradable;
  const buttonTitle = !isTradable ? 'Not Tradable' : isMarketClosed ? 'Market Closed' : undefined;

  return (
    <td className="px-6 py-4 text-right">
      <div className="flex items-center justify-end gap-2">
        <Button
          data-testid="sell-all-button"
          variant="destructive"
          size="sm"
          onClick={() => onSellClick(item)}
          title={buttonTitle || 'Sell All'}
          className="h-8 px-3 text-xs"
          disabled={isActionDisabled}
        >
          <IconTrash className="mr-1 h-3 w-3" />
          Sell All
        </Button>
        <Button
          data-testid="trade-button"
          variant="default"
          size="sm"
          onClick={() => onTradeClick(item.stock_symbol)}
          className="h-8 px-3 text-xs"
          disabled={isActionDisabled}
          title={buttonTitle || 'Trade'}
        >
          Trade
          <IconArrowRight className="ml-1 h-3 w-3" />
        </Button>
      </div>
    </td>
  );
};
