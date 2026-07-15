import { Button } from '@/components/shared/Button';
import { TradeAction } from '@/types';
import styles from './TradePanel.module.css';

export interface TradeButtonProps {
  type: TradeAction;
  onClick?: () => void;
  disabled?: boolean;
}

export const TradeButton = ({ type, onClick, disabled }: TradeButtonProps) => {
  const isBuy = type === TradeAction.BUY;

  return (
    <Button
      onClick={onClick}
      disabled={disabled}
      variant={isBuy ? 'success' : 'destructive'}
      size="lg"
      className={styles.tradeBtn}
    >
      {isBuy ? 'Buy' : 'Sell'}
    </Button>
  );
};
