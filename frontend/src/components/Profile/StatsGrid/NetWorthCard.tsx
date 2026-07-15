import { IconWallet } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import styles from './NetWorthCard.module.css';

interface NetWorthCardProps {
  totalNetWorth: number;
  cash: number;
  assets: number;
}

export const NetWorthCard = ({ totalNetWorth, cash, assets }: NetWorthCardProps) => (
  <Card className={styles.card}>
    <div className={styles.iconWrapper}>
      <IconWallet />
    </div>
    <span className={styles.label}>
      Total Net Worth
    </span>
    <div className={styles.value}>
      ${totalNetWorth.toFixed(2)}
    </div>
    <div className={styles.detailsWrapper}>
      <span className={styles.detailBadge}>
        Cash: <span>${cash.toFixed(2)}</span>
      </span>
      <span className={styles.detailBadge}>
        Assets: <span>${assets.toFixed(2)}</span>
      </span>
    </div>
  </Card>
);
