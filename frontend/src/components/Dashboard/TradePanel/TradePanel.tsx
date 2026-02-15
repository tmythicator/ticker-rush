import { useState } from 'react';
import { Card } from '@/components/ui/card';
import { useTrade } from '@/hooks/useTrade';
import { TradeAction } from '@/types';
import { TradeFooter } from '@/components/Dashboard/TradePanel/TradeFooter';
import { TradeOrderInput } from '@/components/Dashboard/TradePanel/TradeOrderInput';
import { parseTicker } from '@/lib/utils';
import { TradePanelHeader } from '@/components/Dashboard/TradePanel/TradePanelHeader';
import { useAuth } from '@/hooks/useAuth';

export interface TradePanelProps {
  symbol: string | null;
  currentPrice?: number;
  onTradeSuccess?: () => void;
}

export const TradePanel = ({ symbol, currentPrice = 0, onTradeSuccess }: TradePanelProps) => {
  const [quantity, setQuantity] = useState<string>('');
  const { user } = useAuth();
  const buyingPower = user?.balance || 0;
  const { source, symbol: displaySymbol } = parseTicker(symbol || '');

  const { executeTrade, isLoading, error } = useTrade({
    symbol: symbol || '',
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
    <Card className="p-6 flex flex-col h-full relative">
      <TradePanelHeader isLoading={isLoading} />
      <TradeOrderInput
        symbol={displaySymbol.toUpperCase()}
        source={source}
        quantity={quantity}
        setQuantity={setQuantity}
        error={error}
        handleTrade={handleTrade}
      />

      <div className="flex-1" />
      <TradeFooter buyingPower={buyingPower} estCost={estCost} />
    </Card>
  );
};
