import { type User } from '@/types';
import { useDashboardStats } from '@/hooks/useDashboardStats';
import { StatCard } from './StatCard';

interface DashboardStatsProps {
  user: User | null;
}

export const DashboardStats = ({ user }: DashboardStatsProps) => {
  const stats = useDashboardStats(user);

  return (
    <div className="no-scrollbar flex snap-x gap-4 overflow-x-auto pb-4 sm:grid sm:grid-cols-3 sm:pb-0 lg:gap-6">
      {stats.map((stat, i) => (
        <div key={i} className="min-w-[240px] snap-center sm:min-w-0">
          <StatCard {...stat} />
        </div>
      ))}
    </div>
  );
};
