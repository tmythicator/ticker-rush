import { Card } from '@/components/shared/Card';
import { IconBriefcase } from '@/components/icons/CustomIcons';
import { usePortfolioTable } from '@/hooks/usePortfolioTable';
import { PortfolioTable } from './PortfolioTable';
import { SellPositionModal } from '@/components/SellPosition';
import type { PortfolioItem } from '@/types';
import styles from './PortfolioTable.module.css';

interface PortfolioHoldingsProps {
  portfolio: Record<string, PortfolioItem>;
  isReadOnly?: boolean;
}

export const PortfolioHoldings = ({ portfolio, isReadOnly = false }: PortfolioHoldingsProps) => {
  const { items, sellModal, handleTrade, handleSellClick, handleCloseSellModal } =
    usePortfolioTable(portfolio);

  return (
    <Card className={styles.holdingsCard}>
      <div className={styles.holdingsHeaderWrapper}>
        <h3 className={styles.holdingsTitle}>
          <IconBriefcase className={styles.holdingsIcon} />
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
