import { Button } from '@/components/shared/Button';
import { calculateMaxBuyQuantity } from '@/lib/utils';
import styles from './TradePanel.module.css';

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
    <div className={styles.percentageSelectorContainer}>
      {PERCENTAGE_PRESETS.map((pct) => (
        <Button
          key={pct}
          onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, pct))}
          disabled={disabled}
          variant="secondary"
          size="sm"
          className={styles.percentageButton}
        >
          {`${pct * 100}%`}
        </Button>
      ))}
      <Button
        onClick={() => onSelect(calculateMaxBuyQuantity(buyingPower, price, 1.0))}
        disabled={disabled}
        variant="secondary"
        size="sm"
        className={styles.percentageButton}
      >
        MAX
      </Button>
    </div>
  );
};
