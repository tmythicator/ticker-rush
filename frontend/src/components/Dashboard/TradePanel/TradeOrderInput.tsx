import { TradeAction } from '../../../types';
import { TradeButtons } from './TradeButtons';

interface TradeOrderInputProps {
  symbol: string;
  quantity: string;
  setQuantity: (quantity: string) => void;
  error: string | null;
  handleTrade: (action: TradeAction) => void;
}

export const TradeOrderInput = ({
  symbol,
  quantity,
  setQuantity,
  error,
  handleTrade,
}: TradeOrderInputProps) => {
  return (
    <div className="space-y-5 flex-1">
      {error && <div className="text-xs text-red-600 font-bold mb-2">{error}</div>}
      <div>
        <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">
          Symbol
        </label>
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
        <label className="block text-xs font-bold text-slate-400 mb-2 uppercase tracking-wider">
          Quantity
        </label>
        <div className="relative">
          <input
            type="number"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            placeholder="0.0"
            min="0"
            step="any"
            className="w-full bg-white border border-slate-200 rounded-lg px-4 py-3 font-mono text-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all shadow-sm placeholder:text-slate-300"
          />
        </div>
      </div>
      <TradeButtons handleTrade={handleTrade} />
    </div>
  );
};
