import { type User } from '@/types';
import { useDashboardStats } from '@/hooks/useDashboardStats';
import { StatCard } from './StatCard';
import styles from './DashboardStats.module.css';

interface DashboardStatsProps {
  user: User | null;
}

export const DashboardStats = ({ user }: DashboardStatsProps) => {
  const stats = useDashboardStats(user);

  return (
    <div className={styles.container}>
      {stats.map((stat, i) => (
        <div key={i} className={styles.itemWrapper}>
          <StatCard {...stat} />
        </div>
      ))}
    </div>
  );
};
