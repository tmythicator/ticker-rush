import { TradeAction } from '@/types';
import { TradeButton } from './TradeButton';

interface TradeButtonsProps {
  handleTrade: (side: TradeAction) => void;
  disabled?: boolean;
}

export const TradeButtons = ({ handleTrade, disabled }: TradeButtonsProps) => {
  return (
    <div className="grid grid-cols-2 gap-3 pt-2">
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
