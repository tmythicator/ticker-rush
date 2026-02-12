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
    executeTrade({ action: TradeAction.SELL, quantity });
  };

  const price = quote?.price || 0;
  const totalValue = quantity * price;

  const isPriceLoading = isOpen && !quote && !sseError;
  const isPriceError = !!sseError;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Sell ${symbol}?`}>
      <div className="space-y-4">
        <p className="text-muted-foreground text-sm">
          Are you sure you want to sell your entire position of{' '}
          <strong className="text-foreground">{symbol}</strong>?
        </p>

        <div className="bg-muted/50 rounded-lg p-4 space-y-2 border border-border">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Quantity</span>
            <span className="font-mono font-bold text-foreground">{quantity}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Current Price</span>
            {isPriceLoading ? (
              <span className="text-muted-foreground animate-pulse">Loading...</span>
            ) : isPriceError ? (
              <span className="text-destructive font-bold">Unavailable</span>
            ) : (
              <span className="font-mono font-bold text-foreground">${price.toFixed(2)}</span>
            )}
          </div>
          <div className="border-t border-border pt-2 flex justify-between text-sm">
            <span className="font-bold text-foreground">Total Value</span>
            {isPriceLoading ? (
              <span className="text-muted-foreground animate-pulse">Loading...</span>
            ) : isPriceError ? (
              <span className="text-destructive font-bold">Unavailable</span>
            ) : (
              <span className="font-mono font-bold text-foreground">${totalValue.toFixed(2)}</span>
            )}
          </div>
        </div>

        {isPriceError && (
          <div className="text-xs text-destructive bg-destructive/10 p-2 rounded-lg border border-destructive/20">
            Failed to fetch current price. You cannot sell at this time.
          </div>
        )}

        <div className="flex gap-3 pt-2">
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2.5 bg-muted text-foreground font-bold rounded-lg text-sm hover:bg-muted/80 transition-all border border-border"
          >
            Cancel
          </button>
          <button
            onClick={handleSellAll}
            disabled={isTradeLoading || isPriceLoading || isPriceError}
            className="flex-1 px-4 py-2.5 bg-destructive text-destructive-foreground font-bold rounded-lg text-sm hover:bg-destructive/90 shadow-sm transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isTradeLoading ? 'Selling...' : 'Confirm Sell All'}
          </button>
        </div>
      </div>
    </Modal>
  );
};
