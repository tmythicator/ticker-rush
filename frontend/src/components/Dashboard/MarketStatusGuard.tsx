import { IconMoon } from '@/components/icons/CustomIcons';
import clsx from 'clsx';
import type { ReactNode } from 'react';
import { GuardState } from './GuardState';
import styles from './Guard.module.css';

interface MarketStatusGuardProps {
  isParticipating?: boolean;
  isMarketClosed?: boolean;
  isLoadingQuotes?: boolean;
  children: ReactNode;
}

export const MarketStatusGuard = ({ isParticipating, isMarketClosed, isLoadingQuotes, children }: MarketStatusGuardProps) => {

  if (!isParticipating) {
    return (
      <GuardState
        testId="participation-required-guard"
        icon={<IconMoon className={styles.icon} />}
        title="Participation Required"
        description="Join the active ladder to unlock trading features and track your performance."
      />
    );
  }

  if (isLoadingQuotes) {
    return (
      <GuardState
        testId="loading-market-guard"
        icon={<div className={styles.spinner} />}
        title="Loading Market Data"
        description="Fetching the latest quotes..."
      />
    );
  }

  if (isMarketClosed) {
    return (
      <GuardState
        testId="market-closed-guard"
        icon={<IconMoon className={clsx(styles.icon, styles.iconPrimary)} />}
        title="Market Closed"
        description={
          <>
            Trading is currently unavailable.
            <br />
            Please come back during market hours.
          </>
        }
      />
    );
  }

  return <>{children}</>;
};
