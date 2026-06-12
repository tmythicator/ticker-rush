import { IconArrowRight, IconTrash } from '@/components/icons/CustomIcons';
import { Button } from '@/components/shared/Button';
import type { PortfolioItem } from '@/types';

interface RowActionsCellProps {
  item: PortfolioItem;
  isMarketClosed: boolean;
  onSellClick: (item: PortfolioItem) => void;
  onTradeClick: (symbol: string) => void;
}

export const RowActionsCell = ({
  item,
  isMarketClosed,
  onSellClick,
  onTradeClick,
}: RowActionsCellProps) => (
  <td className="px-6 py-4 text-right">
    <div className="flex items-center justify-end gap-2">
      <Button
        variant="destructive"
        size="sm"
        onClick={() => onSellClick(item)}
        title={isMarketClosed ? 'Market Closed' : 'Sell All'}
        className="h-8 px-3 text-xs"
        disabled={isMarketClosed}
      >
        <IconTrash className="w-3 h-3 mr-1" />
        Sell All
      </Button>
      <Button
        variant="default"
        size="sm"
        onClick={() => onTradeClick(item.stock_symbol)}
        className="h-8 px-3 text-xs"
        disabled={isMarketClosed}
        title={isMarketClosed ? 'Market Closed' : 'Trade'}
      >
        Trade
        <IconArrowRight className="w-3 h-3 ml-1" />
      </Button>
    </div>
  </td>
);
