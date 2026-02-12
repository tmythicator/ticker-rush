import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { PortfolioItem } from '../lib/api';
import { IconArrowRight, IconBriefcase, IconTrash } from './icons/CustomIcons';
import { SellPositionModal } from './SellPositionModal';

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
    <div className="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
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
            {Object.values(portfolio).map((item) => (
              <tr key={item.stock_symbol} className="hover:bg-muted/50 transition-colors">
                <td className="px-6 py-4 font-bold text-foreground">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
                      {item.stock_symbol[0]}
                    </div>
                    {item.stock_symbol}
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
                    <button
                      onClick={() => handleSellClick(item)}
                      className="text-xs font-bold text-destructive hover:bg-destructive/10 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1"
                      title="Sell All"
                    >
                      <IconTrash className="w-3 h-3" />
                      Sell All
                    </button>
                    <button
                      onClick={() => handleTrade(item.stock_symbol)}
                      className="text-xs font-bold text-primary hover:bg-primary/10 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1"
                    >
                      Trade
                      <IconArrowRight className="w-3 h-3" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
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
    </div>
  );
};
