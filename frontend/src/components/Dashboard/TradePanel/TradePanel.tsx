import { TradeFooter } from '@/components/Dashboard/TradePanel/TradeFooter';
import { TradeOrderInput } from '@/components/Dashboard/TradePanel/TradeOrderInput';
import { TradePanelHeader } from '@/components/Dashboard/TradePanel/TradePanelHeader';
import { Card } from '@/components/ui/card';
import { useAuth } from '@/hooks/useAuth';
import { useTrade } from '@/hooks/useTrade';
import { parseTicker } from '@/lib/utils';
import { TradeAction, type Quote, type TickerSource } from '@/types';
import { useState } from 'react';

export interface TradePanelProps {
  quote: Quote | null;
  onTradeSuccess?: () => void;
}

export const TradePanel = ({ quote, onTradeSuccess }: TradePanelProps) => {
  const [quantity, setQuantity] = useState<string>('');
  const { user } = useAuth();
  const buyingPower = user?.balance || 0;

  const currentPrice = quote?.price || 0;
  const symbol = quote?.symbol || '';

  // Calculate current position quantity
  const position = user?.portfolio?.[symbol];
  const positionQuantity = position?.quantity || 0;

  const parsed = parseTicker(symbol);
  const source = (quote?.source as TickerSource) || parsed.source;
  const displaySymbol = parsed.symbol;

  const { executeTrade, isLoading, error } = useTrade({
    symbol: symbol,
    onSuccess: () => {
      setQuantity('');
      if (onTradeSuccess) onTradeSuccess();
    },
  });

  const qty = parseFloat(quantity) || 0;
  const estCost = qty * currentPrice;

  const handleTrade = (action: TradeAction) => {
    if (!symbol) return;
    executeTrade({ action, quantity: qty });
  };

  if (!symbol) {
    return (
      <Card className="p-6 flex flex-col h-full items-center justify-center text-muted-foreground">
        No active ticker selected.
      </Card>
    );
  }

  return (
    <Card className="p-6 flex flex-col h-full relative overflow-hidden">
      <TradePanelHeader isLoading={isLoading} />
      <TradeOrderInput
        symbol={displaySymbol.toUpperCase()}
        source={source}
        quantity={quantity}
        setQuantity={setQuantity}
        error={error}
        handleTrade={handleTrade}
        buyingPower={buyingPower}
        price={currentPrice}
        disabled={quote?.is_closed}
        positionQuantity={positionQuantity}
      />

      <div className="flex-1" />
      <TradeFooter buyingPower={buyingPower} estCost={estCost} />
    </Card>
  );
};
