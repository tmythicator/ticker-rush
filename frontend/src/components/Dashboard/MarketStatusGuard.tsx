import { IconMoon } from '@/components/icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import type { Quote, User } from '@/types';
import type { ReactNode } from 'react';
import styles from './MarketStatusGuard.module.css';

interface MarketStatusGuardProps {
  user: User | null;
  quote: Quote | null;
  children: ReactNode;
}

export const MarketStatusGuard = ({ user, quote, children }: MarketStatusGuardProps) => {
  if (!user) return null;

  if (!user.is_participating) {
    return (
      <Card data-testid="participation-required-guard" className={styles.guardCard}>
        <IconMoon className={styles.icon} />
        <h3 className={styles.title}>Participation Required</h3>
        <p className={styles.description}>
          Join the active ladder to unlock trading features and track your performance.
        </p>
      </Card>
    );
  }

  if (quote?.is_closed) {
    return (
      <Card data-testid="market-closed-guard" className={styles.guardCard}>
        <IconMoon className={`${styles.icon} ${styles.iconPrimary}`} />
        <h3 className={styles.title}>Market Closed</h3>
        <p className={styles.description}>
          Trading is currently unavailable.
          <br />
          Please come back during market hours.
        </p>
      </Card>
    );
  }

  if (!quote) {
    return (
      <Card data-testid="loading-market-guard" className={styles.guardCard}>
        <div className={styles.spinner} />
        <h3 className={styles.title}>Loading Market Data</h3>
        <p className={styles.description}>Fetching the latest quotes...</p>
      </Card>
    );
  }

  return <>{children}</>;
};
