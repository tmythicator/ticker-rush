import { usePortfolioValue } from '@/hooks/usePortfolioValue';
import { INITIAL_BALANCE } from '@/lib/constants';
import { calculateInvestedCapital } from '@/lib/utils';
import { type User } from '@/types';
import { NetWorthCard, PortfolioItemsCard, TotalGainLossCard } from './StatsGrid/index';

export const StatsGrid = (user: User) => {
  const portfolioItems = Object.values(user.portfolio ?? {});
  const investedCapital = calculateInvestedCapital(user.portfolio);

  const {
    totalValue: currentPortfolioValue,
    isLoading,
    isError,
  } = usePortfolioValue(user.portfolio);

  const hasItems = portfolioItems.length > 0;
  const isValueSuspect = currentPortfolioValue === 0 && hasItems;
  const shouldUseFallback = isLoading || isError || isValueSuspect;

  const effectivePortfolioValue = shouldUseFallback ? investedCapital : currentPortfolioValue;
  const totalNetWorth = user.balance + effectivePortfolioValue;
  const totalPnL = totalNetWorth - INITIAL_BALANCE;

  return (
    <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
      <NetWorthCard totalNetWorth={totalNetWorth} cash={user.balance} assets={investedCapital} />
      <PortfolioItemsCard count={portfolioItems.length} />
      <TotalGainLossCard totalPnL={totalPnL} />
    </div>
  );
};
