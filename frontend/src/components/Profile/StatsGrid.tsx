import { IconBriefcase, IconTrending, IconWallet } from '@icons/CustomIcons';

import { usePortfolioValue } from '@/hooks/usePortfolioValue';
import { INITIAL_BALANCE } from '@/lib/constants';
import { calculateInvestedCapital, formatCurrencyWithSign } from '@/lib/utils';
import { type User } from '@/types';

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
  const isPnLPositive = totalPnL >= 0;

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div className="bg-gradient-to-br from-primary/10 to-primary/5 border border-primary/20 p-6 rounded-3xl shadow-sm relative overflow-hidden group">
        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <IconWallet className="w-24 h-24 text-primary" />
        </div>
        <div className="text-muted-foreground text-sm font-medium uppercase tracking-wider mb-2">
          Total Net Worth
        </div>
        <div className="text-4xl font-bold font-mono tracking-tight text-foreground">
          ${totalNetWorth.toFixed(2)}
        </div>
        <div className="mt-4 flex flex-wrap items-center gap-2 text-sm">
          <span className="bg-background/50 backdrop-blur-sm px-3 py-1.5 rounded-full border border-border/50 text-muted-foreground">
            Cash: <span className="text-foreground font-medium">${user.balance.toFixed(2)}</span>
          </span>
          <span className="bg-background/50 backdrop-blur-sm px-3 py-1.5 rounded-full border border-border/50 text-muted-foreground">
            Assets:{' '}
            <span className="text-foreground font-medium">${investedCapital.toFixed(2)}</span>
          </span>
        </div>
      </div>

      <div className="bg-card/50 backdrop-blur-sm border border-border/50 p-6 rounded-3xl shadow-sm">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-primary/10 text-primary rounded-2xl">
            <IconBriefcase className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground font-medium">Portfolio Items</div>
            <div className="text-2xl font-bold text-foreground">{portfolioItems.length}</div>
          </div>
        </div>
        <div className="text-sm text-muted-foreground">Active positions in your portfolio.</div>
      </div>

      <div className="bg-card/50 backdrop-blur-sm border border-border/50 p-6 rounded-3xl shadow-sm">
        <div className="flex items-center gap-3 mb-4">
          <div
            className={`p-2 rounded-2xl ${isPnLPositive ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500'}`}
          >
            <IconTrending className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground font-medium">Total Gain/Loss</div>
            <div
              className={`text-2xl font-bold ${isPnLPositive ? 'text-green-500' : 'text-red-500'}`}
            >
              {formatCurrencyWithSign(totalPnL)}
            </div>
          </div>
        </div>
        <div className="text-sm text-muted-foreground">
          Real-time P&L based on current market prices.
        </div>
      </div>
    </div>
  );
};
