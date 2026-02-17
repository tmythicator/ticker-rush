import React from 'react';
import { calculateMaxBuyQuantity } from '@/lib/utils';

interface BuyingPowerControlProps {
  buyingPower: number;
  price: number;
  onSelect: (quantity: string) => void;
  disabled?: boolean;
}

export const MaxButton: React.FC<BuyingPowerControlProps> = ({
  buyingPower,
  price,
  onSelect,
  disabled,
}) => {
  if (!buyingPower || !price || price <= 0) return null;

  return (
    <button
      onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, 1.0))}
      disabled={disabled}
      className="text-xs font-bold text-primary hover:text-primary/80 bg-primary/10 hover:bg-primary/20 px-3 py-1.5 rounded transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none"
    >
      MAX
    </button>
  );
};

export const PercentageSelector: React.FC<BuyingPowerControlProps> = ({
  buyingPower,
  price,
  onSelect,
  disabled,
}) => {
  if (!buyingPower || !price || price <= 0) return null;

  return (
    <div className="flex gap-2 mt-2">
      {[0.1, 0.25, 0.5].map((pct) => (
        <button
          key={pct}
          onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, pct))}
          disabled={disabled}
          className="flex-1 py-1 text-xs font-medium bg-muted hover:bg-muted/80 text-muted-foreground hover:text-foreground rounded transition-colors disabled:opacity-50 disabled:pointer-events-none"
        >
          {(pct * 100).toFixed(0)}%
        </button>
      ))}
    </div>
  );
};
