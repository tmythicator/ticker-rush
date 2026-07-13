import { usePortfolioValue } from '@/hooks/usePortfolioValue';
import { INITIAL_BALANCE } from '@/lib/constants';
import { calculateInvestedCapital } from '@/lib/utils';
import { type PortfolioItem } from '@/types';
import { NetWorthCard, PortfolioItemsCard, TotalGainLossCard } from './StatsGrid/index';

interface StatsGridProps {
  balance: number;
  portfolio?: { [key: string]: PortfolioItem };
}

export const StatsGrid = ({ balance, portfolio = {} }: StatsGridProps) => {
  const portfolioItems = Object.values(portfolio);
  const investedCapital = calculateInvestedCapital(portfolio);

  const { totalValue: currentPortfolioValue, isLoading, isError } = usePortfolioValue(portfolio);

  const hasItems = portfolioItems.length > 0;
  const isValueSuspect = currentPortfolioValue === 0 && hasItems;
  const shouldUseFallback = isLoading || isError || isValueSuspect;

  const effectivePortfolioValue = shouldUseFallback ? investedCapital : currentPortfolioValue;
  const totalNetWorth = balance + effectivePortfolioValue;
  const totalPnL = totalNetWorth - INITIAL_BALANCE;

  return (
    <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
      <NetWorthCard totalNetWorth={totalNetWorth} cash={balance} assets={investedCapital} />
      <PortfolioItemsCard count={portfolioItems.length} />
      <TotalGainLossCard totalPnL={totalPnL} />
    </div>
  );
};
