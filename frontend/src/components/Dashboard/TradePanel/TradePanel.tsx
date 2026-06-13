import { TradeFooter } from '@/components/Dashboard/TradePanel/TradeFooter';
import { TradeOrderInput } from '@/components/Dashboard/TradePanel/TradeOrderInput';
import { TradePanelHeader } from '@/components/Dashboard/TradePanel/TradePanelHeader';
import { Card } from '@/components/shared/Card';
import { useTradePanel } from '@/hooks/useTradePanel';
import type { Quote } from '@/types';

export interface TradePanelProps {
  quote: Quote | null;
  onTradeSuccess?: () => void;
}

export const TradePanel = ({ quote, onTradeSuccess }: TradePanelProps) => {
  const { form, asset, estCost, isLoading, handleTrade, symbol } = useTradePanel({
    quote,
    onTradeSuccess,
  });

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
      <TradeOrderInput asset={asset} form={form} onTrade={handleTrade} />

      <div className="flex-1" />
      <TradeFooter buyingPower={asset.buyingPower || 0} estCost={estCost} />
    </Card>
  );
};
