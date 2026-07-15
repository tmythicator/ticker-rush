import { Button } from '@/components/shared/Button';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { PercentageSelector } from './PercentageSelector';
import styles from './TradePanel.module.css';

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
      <Label className={styles.labelRow}>
        <span>Quantity</span>
        {positionQuantity > 0 && (
          <Button
            onClick={() => setQuantity(positionQuantity.toString())}
            variant="link"
            className={styles.sellAllButton}
          >
            Sell All ({positionQuantity})
          </Button>
        )}
      </Label>
      <div className={styles.quantityInputContainer}>
        <Input
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
