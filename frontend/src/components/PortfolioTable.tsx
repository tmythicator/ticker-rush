import { Briefcase, ArrowRight, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import type { PortfolioItem } from '../lib/api';
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
    <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      <div className="px-6 py-5 border-b border-slate-100 flex items-center justify-between">
        <h3 className="font-bold text-lg text-slate-800 flex items-center gap-2">
          <Briefcase className="w-5 h-5 text-slate-400" />
          Current Holdings
        </h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            <tr className="bg-slate-50 text-slate-500 text-xs uppercase tracking-wider font-bold">
              <th className="px-6 py-4">Asset</th>
              <th className="px-6 py-4 text-right">Quantity</th>
              <th className="px-6 py-4 text-right">Avg Price</th>
              <th className="px-6 py-4 text-right">Total Cost</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-50">
            {Object.values(portfolio).map((item) => (
              <tr key={item.stock_symbol} className="hover:bg-slate-50/50 transition-colors">
                <td className="px-6 py-4 font-bold text-slate-900">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center font-bold text-xs text-slate-600">
                      {item.stock_symbol[0]}
                    </div>
                    {item.stock_symbol}
                  </div>
                </td>
                <td className="px-6 py-4 text-right font-mono text-slate-700">{item.quantity}</td>
                <td className="px-6 py-4 text-right font-mono text-slate-700">
                  ${item.average_price.toFixed(2)}
                </td>
                <td className="px-6 py-4 text-right font-mono text-slate-900 font-bold">
                  ${(item.quantity * item.average_price).toFixed(2)}
                </td>
                <td className="px-6 py-4 text-right">
                  <div className="flex items-center justify-end gap-2">
                    <button
                      onClick={() => handleSellClick(item)}
                      className="text-xs font-bold text-red-600 hover:text-red-700 hover:bg-red-50 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1"
                      title="Sell All"
                    >
                      <Trash2 className="w-3 h-3" />
                      Sell All
                    </button>
                    <button
                      onClick={() => handleTrade(item.stock_symbol)}
                      className="text-xs font-bold text-blue-600 hover:text-blue-700 bg-blue-50 hover:bg-blue-100 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1"
                    >
                      Trade
                      <ArrowRight className="w-3 h-3" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
            {Object.keys(portfolio).length === 0 && (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-slate-400 italic">
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
