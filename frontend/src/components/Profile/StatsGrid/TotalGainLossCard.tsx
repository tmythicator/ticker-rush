import { IconTrending } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import { formatCurrencyWithSign } from '@/lib/utils';
import styles from './TotalGainLossCard.module.css';

interface TotalGainLossCardProps {
  totalPnL: number;
}

export const TotalGainLossCard = ({ totalPnL }: TotalGainLossCardProps) => {
  const isPnLPositive = totalPnL >= 0;
  const trend = isPnLPositive ? 'up' : 'down';

  return (
    <Card className={styles.card}>
      <div className={styles.headerGroup}>
        <div className={styles.iconWrapper} data-trend={trend}>
          <IconTrending />
        </div>
        <div>
          <span className={styles.label}>
            Total Gain/Loss
          </span>
          <div className={styles.value} data-trend={trend}>
            {formatCurrencyWithSign(totalPnL)}
          </div>
        </div>
      </div>
      <p className={styles.description}>
        Real-time P&L based on current market prices.
      </p>
    </Card>
  );
};
