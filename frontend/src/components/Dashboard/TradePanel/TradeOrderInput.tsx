import { TradeAction, type TickerSource } from '@/types';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { TradeButtons } from './TradeButtons';

interface TradeOrderInputProps {
  symbol: string;
  source?: TickerSource;
  quantity: string;
  setQuantity: (quantity: string) => void;
  error: string | null;
  handleTrade: (action: TradeAction) => void;
}

export const TradeOrderInput = ({
  symbol,
  source,
  quantity,
  setQuantity,
  error,
  handleTrade,
}: TradeOrderInputProps) => {
  return (
    <div className="space-y-5 flex-1">
      {error && <div className="text-xs text-red-600 font-bold mb-2">{error}</div>}
      <div>
        <label className="block text-xs font-bold text-muted-foreground mb-2 uppercase tracking-wider">
          Symbol
        </label>
        <div className="w-full bg-muted border border-border rounded-lg px-3 py-3 flex items-center gap-3 opacity-70">
          {source && <SourceBadge source={source} />}
          <input
            type="text"
            value={symbol}
            disabled
            className="flex-1 bg-transparent border-none p-0 font-mono text-sm font-bold text-muted-foreground focus:outline-none"
          />
        </div>
      </div>

      <div>
        <label className="block text-xs font-bold text-muted-foreground mb-2 uppercase tracking-wider">
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
            className="w-full bg-background border border-border rounded-lg px-4 py-3 font-mono text-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all shadow-sm placeholder:text-muted-foreground text-foreground"
          />
        </div>
      </div>
      <TradeButtons handleTrade={handleTrade} />
    </div>
  );
};
