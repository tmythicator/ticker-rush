import { DollarSign, Briefcase, Wallet } from 'lucide-react';
import { type User } from '../../lib/api';
import { calculateInvestedCapital } from '../../lib/utils';
import { StatCard } from './StatCard';

interface DashboardStatsProps {
  user: User | null;
}

export const DashboardStats = ({ user }: DashboardStatsProps) => {
  const portfolioItems = user ? Object.values(user.portfolio ?? {}) : [];
  const portfolioCount = portfolioItems.length;
  const investedCapital = calculateInvestedCapital(user?.portfolio);

  const stats = [
    {
      label: 'Cash Balance',
      value: user ? `$${user.balance.toFixed(2)}` : '--',
      icon: Wallet,
    },
    {
      label: 'Invested Capital',
      value: `$${investedCapital.toFixed(2)}`,
      icon: DollarSign,
    },
    {
      label: 'Open Positions',
      value: portfolioCount.toString(),
      icon: Briefcase,
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
      {stats.map((stat, i) => (
        <StatCard key={i} {...stat} />
      ))}
    </div>
  );
};
