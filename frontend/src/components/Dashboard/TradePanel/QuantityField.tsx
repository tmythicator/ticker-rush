import { Button } from '@/components/shared/Button';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { MaxButton } from './MaxButton';
import { PercentageSelector } from './PercentageSelector';

interface QuantityFieldProps {
  quantity: string;
  setQuantity: (quantity: string) => void;
  buyingPower?: number;
  price?: number;
  disabled?: boolean;
  positionQuantity?: number;
}

export const QuantityField = ({
  quantity,
  setQuantity,
  buyingPower,
  price,
  disabled,
  positionQuantity = 0,
}: QuantityFieldProps) => {
  const showControls = buyingPower !== undefined && price !== undefined && price > 0;

  return (
    <div>
      <Label className="mb-2 block flex items-center justify-between text-xs uppercase tracking-wider text-muted-foreground">
        <span>Quantity</span>
        {positionQuantity > 0 && (
          <Button
            onClick={() => setQuantity(positionQuantity.toString())}
            variant="link"
            className="h-auto p-0 text-xs font-bold text-blue-500 hover:text-blue-600"
          >
            Sell All ({positionQuantity})
          </Button>
        )}
      </Label>
      <div className="relative">
        <Input
          type="number"
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          placeholder="0.0"
          min="0"
          step="any"
          disabled={disabled}
          className="font-mono text-lg [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
        />
        {showControls && (
          <div className="absolute right-3 top-1/2 flex -translate-y-1/2 items-center">
            <MaxButton
              buyingPower={buyingPower}
              price={price}
              onSelect={setQuantity}
              disabled={disabled}
            />
          </div>
        )}
      </div>
      {showControls && (
        <PercentageSelector
          buyingPower={buyingPower}
          price={price}
          onSelect={setQuantity}
          disabled={disabled}
        />
      )}
    </div>
  );
};
