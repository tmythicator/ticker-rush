import { useQuotesSSE } from '../hooks/useQuotesSSE';
import { useTrade } from '../hooks/useTrade';
import { TradeAction } from '../types';
import { Modal } from './Modal';

interface SellPositionModalProps {
  isOpen: boolean;
  onClose: () => void;
  symbol: string;
  quantity: number;
  onSuccess?: () => void;
}

export const SellPositionModal = ({
  isOpen,
  onClose,
  symbol,
  quantity,
  onSuccess,
}: SellPositionModalProps) => {
  const { quote, error: sseError } = useQuotesSSE(isOpen ? symbol : '');

  const { executeTrade, isLoading: isTradeLoading } = useTrade({
    symbol,
    onSuccess: () => {
      if (onSuccess) onSuccess();
      onClose();
    },
  });

  const handleSellAll = () => {
    executeTrade(TradeAction.SELL, quantity);
  };

  const price = quote?.price || 0;
  const totalValue = quantity * price;

  const isPriceLoading = isOpen && !quote && !sseError;
  const isPriceError = !!sseError;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Sell ${symbol}?`}>
      <div className="space-y-4">
        <p className="text-slate-600 text-sm">
          Are you sure you want to sell your entire position of{' '}
          <strong className="text-slate-900">{symbol}</strong>?
        </p>

        <div className="bg-slate-50 rounded-lg p-4 space-y-2 border border-slate-100">
          <div className="flex justify-between text-sm">
            <span className="text-slate-500">Quantity</span>
            <span className="font-mono font-bold text-slate-700">{quantity}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-slate-500">Current Price</span>
            {isPriceLoading ? (
              <span className="text-slate-400 animate-pulse">Loading...</span>
            ) : isPriceError ? (
              <span className="text-red-500 font-bold">Unavailable</span>
            ) : (
              <span className="font-mono font-bold text-slate-700">${price.toFixed(2)}</span>
            )}
          </div>
          <div className="border-t border-slate-200 pt-2 flex justify-between text-sm">
            <span className="font-bold text-slate-900">Total Value</span>
            {isPriceLoading ? (
              <span className="text-slate-400 animate-pulse">Loading...</span>
            ) : isPriceError ? (
              <span className="text-red-500 font-bold">Unavailable</span>
            ) : (
              <span className="font-mono font-bold text-slate-900">${totalValue.toFixed(2)}</span>
            )}
          </div>
        </div>

        {isPriceError && (
          <div className="text-xs text-red-600 bg-red-50 p-2 rounded border border-red-100">
            Failed to fetch current price. You cannot sell at this time.
          </div>
        )}

        <div className="flex gap-3 pt-2">
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2.5 bg-white border border-slate-200 text-slate-700 font-bold rounded-lg text-sm hover:bg-slate-50 transition-all"
          >
            Cancel
          </button>
          <button
            onClick={handleSellAll}
            disabled={isTradeLoading || isPriceLoading || isPriceError}
            className="flex-1 px-4 py-2.5 bg-red-500 text-white font-bold rounded-lg text-sm hover:bg-red-600 shadow-sm shadow-red-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isTradeLoading ? 'Selling...' : 'Confirm Sell All'}
          </button>
        </div>
      </div>
    </Modal>
  );
};
