import { TradeAction } from '../../../types';
import { TradeButton } from './TradeButton';

interface TradeButtonsProps {
  handleTrade: (side: TradeAction) => void;
}

export const TradeButtons = ({ handleTrade }: TradeButtonsProps) => {
  return (
    <div className="pt-2 grid grid-cols-2 gap-3">
      <TradeButton type={TradeAction.BUY} onClick={() => handleTrade(TradeAction.BUY)} />
      <TradeButton type={TradeAction.SELL} onClick={() => handleTrade(TradeAction.SELL)} />
    </div>
  );
};
