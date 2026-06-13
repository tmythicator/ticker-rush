import { type User } from '@/types';
import { useDashboardStats } from '@/hooks/useDashboardStats';
import { StatCard } from './StatCard';

interface DashboardStatsProps {
  user: User | null;
}

export const DashboardStats = ({ user }: DashboardStatsProps) => {
  const stats = useDashboardStats(user);

  return (
    <div className="flex sm:grid sm:grid-cols-3 gap-4 lg:gap-6 overflow-x-auto pb-4 sm:pb-0 snap-x no-scrollbar">
      {stats.map((stat, i) => (
        <div key={i} className="snap-center min-w-[240px] sm:min-w-0">
          <StatCard {...stat} />
        </div>
      ))}
    </div>
  );
};
