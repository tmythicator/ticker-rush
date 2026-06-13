import { Card } from '@/components/shared/Card';
import { IconBriefcase } from '@/components/icons/CustomIcons';
import { usePortfolioTable } from '@/hooks/usePortfolioTable';
import { PortfolioTable } from './PortfolioTable';
import { SellPositionModal } from '@/components/SellPosition';
import type { PortfolioItem } from '@/types';

interface PortfolioHoldingsProps {
  portfolio: Record<string, PortfolioItem>;
  isReadOnly?: boolean;
}

export const PortfolioHoldings = ({ portfolio, isReadOnly = false }: PortfolioHoldingsProps) => {
  const { items, sellModal, handleTrade, handleSellClick, handleCloseSellModal } =
    usePortfolioTable(portfolio);

  return (
    <Card className="overflow-hidden">
      <div className="px-6 py-5 border-b border-border flex items-center justify-between">
        <h3 className="font-bold text-lg text-foreground flex items-center gap-2">
          <IconBriefcase className="w-5 h-5 text-muted-foreground" />
          Current Holdings
        </h3>
      </div>

      <PortfolioTable
        items={items}
        isReadOnly={isReadOnly}
        onSellClick={handleSellClick}
        onTradeClick={handleTrade}
      />

      {sellModal.item && (
        <SellPositionModal
          isOpen={sellModal.isOpen}
          onClose={handleCloseSellModal}
          symbol={sellModal.item.stock_symbol}
          quantity={sellModal.item.quantity}
          onSuccess={handleCloseSellModal}
        />
      )}
    </Card>
  );
};
