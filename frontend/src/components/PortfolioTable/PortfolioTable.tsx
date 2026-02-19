import { PortfolioRow } from '@/components/PortfolioTable/PortfolioRow';
import { SellPositionModal } from '@/components/SellPositionModal';
import { IconBriefcase } from '@/components/icons/CustomIcons';
import { Card } from '@/components/ui/card';
import type { PortfolioItem } from '@/types';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface PortfolioTableProps {
  portfolio: Record<string, PortfolioItem>;
  isReadOnly?: boolean;
}

export const PortfolioTable = ({ portfolio, isReadOnly = false }: PortfolioTableProps) => {
  const navigate = useNavigate();
  const [sellModal, setSellModal] = useState<{ isOpen: boolean; item?: PortfolioItem }>({
    isOpen: false,
  });

  const handleTrade = (symbol: string) => {
    navigate(`/trade?symbol=${symbol}`);
  };

  const handleSellClick = (item: PortfolioItem) => {
    setSellModal({ isOpen: true, item });
  };

  return (
    <Card className="overflow-hidden">
      <div className="px-6 py-5 border-b border-border flex items-center justify-between">
        <h3 className="font-bold text-lg text-foreground flex items-center gap-2">
          <IconBriefcase className="w-5 h-5 text-muted-foreground" />
          Current Holdings
        </h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            <tr className="bg-muted text-muted-foreground text-xs uppercase tracking-wider font-bold">
              <th className="px-6 py-4">Asset</th>
              <th className="px-6 py-4 text-right">Quantity</th>
              <th className="px-6 py-4 text-right">Avg Price</th>
              <th className="px-6 py-4 text-right">Current Price</th>
              <th className="px-6 py-4 text-right">Market Value</th>
              <th className="px-6 py-4 text-right">P&L</th>
              {!isReadOnly && <th className="px-6 py-4 text-right">Actions</th>}
            </tr>
          </thead>

          <tbody className="divide-y divide-border">
            {Object.values(portfolio).map((item) => (
              <PortfolioRow
                key={item.stock_symbol}
                item={item}
                onSellClick={handleSellClick}
                onTradeClick={handleTrade}
                isReadOnly={isReadOnly}
              />
            ))}
            {Object.keys(portfolio).length === 0 && (
              <tr>
                <td
                  colSpan={isReadOnly ? 6 : 7}
                  className="px-6 py-12 text-center text-muted-foreground italic"
                >
                  No assets found in your portfolio.{!isReadOnly && ' Start trading!'}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {sellModal.item && (
        <SellPositionModal
          isOpen={sellModal.isOpen}
          onClose={() => setSellModal({ isOpen: false })}
          symbol={sellModal.item.stock_symbol}
          quantity={sellModal.item.quantity}
          onSuccess={() => setSellModal({ isOpen: false })}
        />
      )}
    </Card>
  );
};
