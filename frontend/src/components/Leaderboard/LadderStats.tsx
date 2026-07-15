import { IconCalendar, IconWallet } from '@/components/icons/CustomIcons';
import { formatLocalTime } from '@/lib/utils';
import styles from './LadderStats.module.css';

interface LadderStatsProps {
  endTime?: Date;
  initialBalance?: number;
}

export const LadderStats = ({ endTime, initialBalance }: LadderStatsProps) => {
  return (
    <div className={styles.statsContainer}>
      <div className={styles.statsCard}>
        <div className={styles.statsIconWrapper} data-type="time">
          <IconCalendar className={styles.statsIcon} />
        </div>
        <div>
          <div className={styles.statsLabel}>Competition Ends</div>
          <div className={styles.statsValue}>
            {endTime ? formatLocalTime(endTime.getTime() / 1000) : 'N/A'}
          </div>
        </div>
      </div>

      <div className={styles.statsCard}>
        <div className={styles.statsIconWrapper} data-type="capital">
          <IconWallet className={styles.statsIcon} />
        </div>
        <div>
          <div className={styles.statsLabel}>Starting Capital</div>
          <div className={styles.statsValue}>${initialBalance?.toLocaleString() ?? '0'}</div>
        </div>
      </div>
    </div>
  );
};
