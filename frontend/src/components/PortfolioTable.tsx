import { IconArrowRight, IconBriefcase, IconTrash } from '@/components/icons/CustomIcons';
import { SellPositionModal } from '@/components/SellPositionModal';
import { parseTicker } from '@/lib/utils';
import type { PortfolioItem } from '@/types';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { SourceBadge } from '@/components/shared/SourceBadge';

interface PortfolioTableProps {
  portfolio: Record<string, PortfolioItem>;
}

export const PortfolioTable = ({ portfolio }: PortfolioTableProps) => {
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
              <th className="px-6 py-4 text-right">Total Cost</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {Object.values(portfolio).map((item) => {
              const { source, symbol } = parseTicker(item.stock_symbol);
              return (
                <tr key={item.stock_symbol} className="hover:bg-muted/50 transition-colors">
                  <td className="px-6 py-4 font-bold text-foreground">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
                        {symbol[0].toUpperCase()}
                      </div>
                      <div className="flex flex-col items-start">
                        <span className="text-sm font-bold">{symbol.toUpperCase()}</span>
                        <SourceBadge source={source} />
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-right font-mono text-muted-foreground">
                    {item.quantity}
                  </td>
                  <td className="px-6 py-4 text-right font-mono text-muted-foreground">
                    ${item.average_price.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 text-right font-mono text-foreground font-bold">
                    ${(item.quantity * item.average_price).toFixed(2)}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleSellClick(item)}
                        title="Sell All"
                        className="h-8 px-3 text-xs"
                      >
                        <IconTrash className="w-3 h-3 mr-1" />
                        Sell All
                      </Button>
                      <Button
                        variant="default"
                        size="sm"
                        onClick={() => handleTrade(item.stock_symbol)}
                        className="h-8 px-3 text-xs"
                      >
                        Trade
                        <IconArrowRight className="w-3 h-3 ml-1" />
                      </Button>
                    </div>
                  </td>
                </tr>
              );
            })}
            {Object.keys(portfolio).length === 0 && (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-muted-foreground italic">
                  No assets found in your portfolio. Start trading!
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
