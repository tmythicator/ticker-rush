import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { useTrade } from '@/hooks/useTrade';
import type { Quote, TickerSource, TradeAction } from '@/types';

interface UseTradePanelProps {
  quote: Quote | null;
  onTradeSuccess?: () => void;
}

export const useTradePanel = ({ quote, onTradeSuccess }: UseTradePanelProps) => {
  const [quantity, setQuantity] = useState<string>('');
  const { user } = useAuth();
  const buyingPower = user?.balance || 0;

  const currentPrice = quote?.price || 0;
  const symbol = quote?.symbol || '';

  // Calculate current position quantity
  const position = user?.portfolio?.[symbol];
  const positionQuantity = position?.quantity || 0;

  const source = (quote?.source || 'Finnhub') as TickerSource;
  const displaySymbol = quote?.symbol || symbol;

  const { executeTrade, isLoading, error } = useTrade({
    symbol: symbol,
    onSuccess: () => {
      setQuantity('');
      if (onTradeSuccess) onTradeSuccess();
    },
  });

  const qty = parseFloat(quantity) || 0;
  const roundedQty = Math.round(qty * 100000000) / 100000000;
  const estCost = roundedQty * currentPrice;

  const handleTrade = (action: TradeAction) => {
    if (!symbol) return;
    executeTrade({ action, quantity: roundedQty });
  };

  return {
    form: {
      quantity,
      setQuantity,
      error,
      disabled: quote?.is_closed,
    },
    asset: {
      symbol: displaySymbol.toUpperCase(),
      source,
      price: currentPrice,
      positionQuantity,
      buyingPower,
    },
    estCost,
    isLoading,
    handleTrade,
    symbol,
  };
};
