import { useState } from 'react';
import { TrendingUp, RefreshCcw } from 'lucide-react';
import { TradeButton } from './TradeButton';
import { buyStock, sellStock } from '../../lib/api';

export interface TradePanelProps {
    userId: number;
    symbol: string;
    currentPrice?: number;
    buyingPower?: number;
    onTradeSuccess?: () => void;
}

export const TradePanel = ({ userId, symbol, currentPrice = 0, buyingPower = 0, onTradeSuccess }: TradePanelProps) => {
    const [quantity, setQuantity] = useState<string>('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const qty = parseInt(quantity) || 0;
    const estCost = qty * currentPrice;

    const handleTrade = async (side: 'BUY' | 'SELL') => {
        if (!qty) return;
        setError(null);
        setIsLoading(true);
        try {
            if (side === 'BUY') {
                await buyStock(userId, symbol, qty);
            } else {
                await sellStock(userId, symbol, qty);
            }
            setQuantity('');
            if (onTradeSuccess) onTradeSuccess();
        } catch (e: any) {
            setError(e.message || 'Trade failed');
        } finally {
            setIsLoading(false);
        }
    };

    const formatCurrency = (val: number) =>
        new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

    return (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 flex flex-col h-full">
            <div className="flex items-center justify-between mb-6">
                <h2 className="font-bold text-slate-800 flex items-center gap-2">
                    <TrendingUp className="w-4 h-4 text-blue-600" />
                    Place Order
                </h2>
                {isLoading && <RefreshCcw className="w-4 h-4 animate-spin text-slate-400" />}
            </div>

            <div className="space-y-5 flex-1">
                {error && <div className="text-xs text-red-600 font-bold mb-2">{error}</div>}

                <div>
                    <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">Symbol</label>
                    <div className="relative">
                        <input
                            type="text"
                            value={symbol}
                            disabled
                            className="w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-3 font-mono text-sm font-bold text-slate-700 opacity-70"
                        />
                        <div className="absolute right-3 top-3 text-xs font-bold text-slate-400">STOCK</div>
                    </div>
                </div>

                <div>
                    <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">Quantity</label>
                    <div className="relative">
                        <input
                            type="number"
                            value={quantity}
                            onChange={(e) => setQuantity(e.target.value)}
                            placeholder="0"
                            min="1"
                            className="w-full bg-white border border-slate-200 rounded-lg px-4 py-3 font-mono text-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all shadow-sm placeholder:text-slate-300"
                        />
                    </div>
                </div>

                <div className="pt-2 grid grid-cols-2 gap-3">
                    <TradeButton type="BUY" onClick={() => handleTrade('BUY')} />
                    <TradeButton type="SELL" onClick={() => handleTrade('SELL')} />
                </div>
            </div>

            <div className="mt-auto pt-6 border-t border-slate-100 space-y-2">
                <div className="flex justify-between text-sm">
                    <span className="text-slate-500">Buying Power</span>
                    <span className="font-mono font-medium text-slate-700">{formatCurrency(buyingPower)}</span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-slate-500">Est. Cost</span>
                    <span className="font-mono font-bold text-slate-900">{formatCurrency(estCost)}</span>
                </div>
            </div>
        </div>
    );
};
