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
    <div className="mt-2.5 flex gap-1.5">
      {PERCENTAGE_PRESETS.map((pct) => (
        <Button
          key={pct}
          onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, pct))}
          disabled={disabled}
          variant="secondary"
          size="sm"
          className="h-7 flex-1 rounded-md border border-border/40 bg-muted/40 text-[10px] font-bold tracking-wider transition-all duration-200 hover:border-primary/30 hover:bg-primary/10 hover:text-primary active:scale-95"
        >
          {`${pct * 100}%`}
        </Button>
      ))}
    </div>
  );
};
