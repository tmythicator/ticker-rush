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
      <div className="flex items-center justify-between border-b border-border px-6 py-5">
        <h3 className="flex items-center gap-2 text-lg font-bold text-foreground">
          <IconBriefcase className="h-5 w-5 text-muted-foreground" />
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
