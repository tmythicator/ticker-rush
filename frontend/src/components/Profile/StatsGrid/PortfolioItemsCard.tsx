import { IconBriefcase } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import styles from './PortfolioItemsCard.module.css';

interface PortfolioItemsCardProps {
  count: number;
}

export const PortfolioItemsCard = ({ count }: PortfolioItemsCardProps) => (
  <Card className={styles.card}>
    <div className={styles.headerGroup}>
      <div className={styles.iconWrapper}>
        <IconBriefcase />
      </div>
      <div>
        <span className={styles.label}>Portfolio Items</span>
        <div className={styles.value}>{count}</div>
      </div>
    </div>
    <p className={styles.description}>Active positions in your portfolio.</p>
  </Card>
);
