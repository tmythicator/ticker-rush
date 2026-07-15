import { IconMoon } from '@/components/icons/CustomIcons';
import type { Quote, User } from '@/types';
import clsx from 'clsx';
import type { ReactNode } from 'react';
import { GuardState } from './GuardState';
import styles from './Guard.module.css';

interface MarketStatusGuardProps {
  user: User | null;
  quote: Quote | null;
  children: ReactNode;
}

export const MarketStatusGuard = ({ user, quote, children }: MarketStatusGuardProps) => {
  if (!user) return null;

  if (!user.is_participating) {
    return (
      <GuardState
        testId="participation-required-guard"
        icon={<IconMoon className={styles.icon} />}
        title="Participation Required"
        description="Join the active ladder to unlock trading features and track your performance."
      />
    );
  }

  if (!quote) {
    return (
      <GuardState
        testId="loading-market-guard"
        icon={<div className={styles.spinner} />}
        title="Loading Market Data"
        description="Fetching the latest quotes..."
      />
    );
  }

  if (quote.is_closed) {
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
