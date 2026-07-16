import { usePortfolioValue } from '@/hooks/usePortfolioValue';
import { INITIAL_BALANCE } from '@/lib/constants';
import { calculateInvestedCapital } from '@/lib/utils';
import { type PortfolioItem } from '@/types';
import { NetWorthCard, PortfolioItemsCard, TotalGainLossCard } from './StatsGrid/index';
import styles from './StatsGrid.module.css';

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
    <section aria-label="Financial Summary" className={styles.statsSection}>
      <ul className={styles.statsContainer}>
        <li className={styles.statItem}>
          <NetWorthCard totalNetWorth={totalNetWorth} cash={balance} assets={investedCapital} />
        </li>
        <li className={styles.statItem}>
          <PortfolioItemsCard count={portfolioItems.length} />
        </li>
        <li className={styles.statItem}>
          <TotalGainLossCard totalPnL={totalPnL} />
        </li>
      </ul>
    </section>
  );
};
