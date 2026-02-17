import { TradeAction, type TickerSource } from '@/types';
import { SourceBadge } from '@/components/shared/SourceBadge';
import { TradeButtons } from './TradeButtons';
import { MaxButton, PercentageSelector } from './BuyingPowerControls';

interface TradeOrderInputProps {
  symbol: string;
  source?: TickerSource;
  quantity: string;
  setQuantity: (quantity: string) => void;
  error: string | null;
  handleTrade: (action: TradeAction) => void;
  buyingPower?: number;
  price?: number;
  disabled?: boolean;
}

export const TradeOrderInput = ({
  symbol,
  source,
  quantity,
  setQuantity,
  error,
  handleTrade,
  buyingPower,
  price,
  disabled,
}: TradeOrderInputProps) => {
  return (
    <div className={`space-y-5 flex-1 ${disabled ? 'opacity-50 pointer-events-none' : ''}`}>
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
            disabled={disabled}
            className="w-full bg-background border border-border rounded-lg px-4 py-3 font-mono text-lg focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all shadow-sm placeholder:text-muted-foreground text-foreground [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none disabled:cursor-not-allowed disabled:bg-muted"
          />
          {buyingPower !== undefined && price && price > 0 && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center">
              <MaxButton
                buyingPower={buyingPower}
                price={price}
                onSelect={setQuantity}
                disabled={disabled}
              />
            </div>
          )}
        </div>
        {buyingPower !== undefined && price && price > 0 && (
          <PercentageSelector
            buyingPower={buyingPower}
            price={price}
            onSelect={setQuantity}
            disabled={disabled}
          />
        )}
      </div>
      <TradeButtons handleTrade={handleTrade} disabled={disabled} />
    </div>
  );
};
