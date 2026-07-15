import { IconArrowRight, IconTrash } from '@/components/icons/CustomIcons';
import { Button } from '@/components/shared/Button';
import type { PortfolioItem } from '@/types';
import styles from './PortfolioTable.module.css';

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
    <td className={styles.cell} data-align="center">
      <div className={styles.actionsContainer}>
        <Button
          data-testid="sell-all-button"
          variant="destructive"
          size="sm"
          onClick={() => onSellClick(item)}
          title={buttonTitle || 'Sell All'}
          className={styles.actionBtn}
          disabled={isActionDisabled}
        >
          <IconTrash className={styles.actionIconLeft} />
          Sell All
        </Button>
        <Button
          data-testid="trade-button"
          variant="default"
          size="sm"
          onClick={() => onTradeClick(item.stock_symbol)}
          className={styles.actionBtn}
          disabled={isActionDisabled}
          title={buttonTitle || 'Trade'}
        >
          Trade
          <IconArrowRight className={styles.actionIconRight} />
        </Button>
      </div>
    </td>
  );
};
