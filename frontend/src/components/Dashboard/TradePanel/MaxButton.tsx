import { Button } from '@/components/shared/Button';
import { calculateMaxBuyQuantity } from '@/lib/utils';

interface MaxButtonProps {
  buyingPower: number;
  price: number;
  onSelect: (quantity: string) => void;
  disabled?: boolean;
}

export const MaxButton = ({ buyingPower, price, onSelect, disabled }: MaxButtonProps) => {
  if (!buyingPower || !price || price <= 0) return null;

  return (
    <Button
      onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, 1.0))}
      disabled={disabled}
      variant="secondary"
      size="sm"
      className="text-xs font-bold text-primary hover:text-primary/80 bg-primary/10 hover:bg-primary/20 h-7 px-3 rounded transition-all active:scale-95"
    >
      MAX
    </Button>
  );
};
