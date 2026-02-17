import { TradeAction } from '@/types';

export interface TradeButtonProps {
  type: TradeAction;
  onClick?: () => void;
  disabled?: boolean;
}

export const TradeButton = ({ type, onClick, disabled }: TradeButtonProps) => {
  const isBuy = type === TradeAction.BUY;

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`w-full font-bold py-3 rounded-lg transition-all shadow-sm active:scale-95 transform duration-100 flex items-center justify-center gap-2 disabled:opacity-50 disabled:pointer-events-none disabled:shadow-none ${
        isBuy
          ? 'bg-green-600 hover:bg-green-700 text-white shadow-green-100'
          : 'bg-red-600 hover:bg-red-700 text-white shadow-red-100'
      }`}
    >
      {type}
    </button>
  );
};
