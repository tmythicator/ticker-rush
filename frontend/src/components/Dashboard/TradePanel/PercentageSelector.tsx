import { Button } from '@/components/shared/Button';
import { calculateMaxBuyQuantity } from '@/lib/utils';

interface PercentageSelectorProps {
  buyingPower: number;
  price: number;
  onSelect: (quantity: string) => void;
  disabled?: boolean;
}

const PERCENTAGE_PRESETS = [0.1, 0.25, 0.5, 0.75];

export const PercentageSelector = ({
  buyingPower,
  price,
  onSelect,
  disabled,
}: PercentageSelectorProps) => {
  if (!buyingPower || !price || price <= 0) return null;

  return (
    <div className="flex gap-1.5 mt-2.5">
      {PERCENTAGE_PRESETS.map((pct) => (
        <Button
          key={pct}
          onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, pct))}
          disabled={disabled}
          variant="secondary"
          size="sm"
          className="flex-1 h-7 text-[10px] font-bold tracking-wider bg-muted/40 hover:bg-primary/10 hover:text-primary border border-border/40 hover:border-primary/30 rounded-md transition-all duration-200 active:scale-95"
        >
          {`${pct * 100}%`}
        </Button>
      ))}
    </div>
  );
};
