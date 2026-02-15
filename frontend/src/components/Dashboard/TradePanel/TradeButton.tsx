import { TradeAction } from '@/types';

export interface TradeButtonProps {
  type: TradeAction;
  onClick?: () => void;
}

export const TradeButton = ({ type, onClick }: TradeButtonProps) => {
  const isBuy = type === TradeAction.BUY;

  return (
    <button
      onClick={onClick}
      className={`w-full font-bold py-3 rounded-lg transition-all shadow-sm active:scale-95 transform duration-100 flex items-center justify-center gap-2 ${
        isBuy
          ? 'bg-green-600 hover:bg-green-700 text-white shadow-green-100'
          : 'bg-red-600 hover:bg-red-700 text-white shadow-red-100'
      }`}
    >
      {type}
    </button>
  );
};
