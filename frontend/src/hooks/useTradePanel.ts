import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { useTrade } from '@/hooks/useTrade';
import type { Quote, TickerSource, TradeAction, User } from '@/types';

interface UseTradePanelProps {
  quote: Quote | null;
  onTradeSuccess?: () => void;
}

export const useTradePanel = ({ quote, onTradeSuccess }: UseTradePanelProps) => {
  const [quantity, setQuantity] = useState<string>('');
  const { user } = useAuth();

  const meta = getTradeMetadata(user, quote);
  const { roundedQty, estCost } = calculateTradeMath(quantity, meta.currentPrice);

  const { executeTrade, isLoading, error } = useTrade({
    symbol: meta.symbol,
    onSuccess: () => {
      setQuantity('');
      onTradeSuccess?.();
    },
  });

  const handleTrade = (action: TradeAction) => {
    if (!meta.symbol) return;
    executeTrade({ action, quantity: roundedQty });
  };

  return {
    form: {
      quantity,
      setQuantity,
      error,
      disabled: meta.isClosed,
    },
    asset: {
      symbol: meta.displaySymbol,
      source: meta.source,
      price: meta.currentPrice,
      positionQuantity: meta.positionQuantity,
      buyingPower: meta.buyingPower,
    },
    estCost,
    isLoading,
    handleTrade,
    symbol: meta.symbol,
  };
};

// eslint-disable-next-line complexity
const getTradeMetadata = (user: User | null | undefined, quote: Quote | null | undefined) => {
  const symbol = quote?.symbol || '';
  const currentPrice = quote?.price || 0;

  return {
    buyingPower: user?.balance || 0,
    currentPrice,
    symbol,
    positionQuantity: user?.portfolio?.[symbol]?.quantity || 0,
    source: (quote?.source || 'Finnhub') as TickerSource,
    displaySymbol: (quote?.symbol || symbol).toUpperCase(),
    isClosed: quote?.is_closed ?? false,
  };
};

const calculateTradeMath = (quantityStr: string, price: number) => {
  const qty = parseFloat(quantityStr) || 0;
  const roundedQty = Math.round(qty * 100000000) / 100000000;
  return {
    roundedQty,
    estCost: roundedQty * price,
  };
};
