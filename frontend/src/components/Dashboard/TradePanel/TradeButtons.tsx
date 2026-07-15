import { TradeAction } from '@/types';
import { TradeButton } from './TradeButton';
import styles from './TradePanel.module.css';

interface TradeButtonsProps {
  handleTrade: (side: TradeAction) => void;
  disabled?: boolean;
}

export const TradeButtons = ({ handleTrade, disabled }: TradeButtonsProps) => {
  return (
    <div className={styles.gridButtons}>
      <TradeButton
        type={TradeAction.BUY}
        onClick={() => handleTrade(TradeAction.BUY)}
        disabled={disabled}
      />
      <TradeButton
        type={TradeAction.SELL}
        onClick={() => handleTrade(TradeAction.SELL)}
        disabled={disabled}
      />
    </div>
  );
};
