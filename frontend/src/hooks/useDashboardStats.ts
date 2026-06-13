import { useMemo } from 'react';
import { calculateInvestedCapital } from '@/lib/utils';
import { IconBriefcase, IconDollarSign, IconWallet } from '@icons/CustomIcons';
import type { User } from '@/types';

export const useDashboardStats = (user: User | null) => {
  return useMemo(() => {
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
};
