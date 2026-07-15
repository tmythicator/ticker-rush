import { TradeFooter } from '@/components/Dashboard/TradePanel/TradeFooter';
import { TradeOrderInput } from '@/components/Dashboard/TradePanel/TradeOrderInput';
import { TradePanelHeader } from '@/components/Dashboard/TradePanel/TradePanelHeader';
import { Card } from '@/components/shared/Card';
import { useTradePanel } from '@/hooks/useTradePanel';
import type { Quote } from '@/types';
import styles from './TradePanel.module.css';

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
      <Card className={styles.emptyState}>
        No active ticker selected.
      </Card>
    );
  }

  return (
    <Card className={styles.panelCard}>
      <TradePanelHeader isLoading={isLoading} />
      <TradeOrderInput asset={asset} form={form} onTrade={handleTrade} />

      <TradeFooter buyingPower={asset.buyingPower || 0} estCost={estCost} />
    </Card>
  );
};
