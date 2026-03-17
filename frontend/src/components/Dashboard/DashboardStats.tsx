import { calculateInvestedCapital } from '@/lib/utils';
import { type User } from '@/types';
import { IconBriefcase, IconDollarSign, IconWallet } from '@icons/CustomIcons';
import { useMemo } from 'react';
import { StatCard } from './StatCard';

interface DashboardStatsProps {
  user: User | null;
}

export const DashboardStats = ({ user }: DashboardStatsProps) => {
  const stats = useMemo(() => {
    const portfolio = user?.portfolio ?? {};
    const portfolioCount = Object.keys(portfolio).length;
    const investedCapital = calculateInvestedCapital(portfolio);

    return [
      {
        label: 'Cash Balance',
        value: user ? `$${user.balance.toFixed(2)}` : '--',
        icon: IconWallet,
      },
      {
        label: 'Invested Capital',
        value: `$${investedCapital.toFixed(2)}`,
        icon: IconDollarSign,
      },
      {
        label: 'Open Positions',
        value: portfolioCount.toString(),
        icon: IconBriefcase,
      },
    ];
  }, [user]);

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
