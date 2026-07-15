import { Button } from '@/components/shared/Button';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { PercentageSelector } from './PercentageSelector';
import styles from './TradePanel.module.css';
import { useId } from 'react';

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
  const inputId = useId();
  const showControls = buyingPower !== undefined && price !== undefined && price > 0;

  return (
    <div>
      <div className={styles.labelRow}>
        <Label htmlFor={inputId}>Quantity</Label>
        {positionQuantity > 0 && (
          <Button
            onClick={() => setQuantity(positionQuantity.toString())}
            variant="link"
            type="button"
            className={styles.sellAllButton}
            aria-label={`Sell all shares of current position ${positionQuantity}`}
          >
            Sell All ({positionQuantity})
          </Button>
        )}
      </div>
      <div className={styles.quantityInputContainer}>
        <Input
          id={inputId}
          type="number"
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          placeholder="0.0"
          min="0"
          step="any"
          disabled={disabled}
          className={styles.quantityInput}
        />
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
