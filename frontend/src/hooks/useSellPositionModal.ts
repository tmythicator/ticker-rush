import { useQuotesSSE } from '@/hooks/useQuotesSSE';
import { useTickers } from '@/hooks/useTickers';
import { useTrade } from '@/hooks/useTrade';
import { TradeAction } from '@/types';

interface UseSellPositionModalProps {
  isOpen: boolean;
  symbol: string;
  quantity: number;
  onClose: () => void;
  onSuccess?: () => void;
}

export const useSellPositionModal = ({
  isOpen,
  symbol,
  quantity,
  onClose,
  onSuccess,
}: UseSellPositionModalProps) => {
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

  const { data: config } = useTickers();
  const tickers = config || [];
  const tickerInfo = tickers.find((t) => t.symbol === symbol);
  const displaySymbol = (tickerInfo?.symbol || symbol).toUpperCase();

  return {
    displaySymbol,
    price,
    totalValue,
    isPriceLoading,
    isPriceError,
    isTradeLoading,
    handleSellAll,
  };
};
