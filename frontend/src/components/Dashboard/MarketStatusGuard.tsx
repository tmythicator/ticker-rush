import { IconMoon } from '@/components/icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import type { Quote, User } from '@/types';
import type { ReactNode } from 'react';

interface MarketStatusGuardProps {
  user: User | null;
  quote: Quote | null;
  children: ReactNode;
}

export const MarketStatusGuard = ({ user, quote, children }: MarketStatusGuardProps) => {
  if (!user) return null;

  if (!user.is_participating) {
    return (
      <Card
        data-testid="participation-required-guard"
        className="flex h-full min-h-[300px] flex-col items-center justify-center p-6 text-center"
      >
        <IconMoon className="mb-4 h-8 w-8 text-secondary/50" />
        <h3 className="text-xl font-bold text-foreground">Participation Required</h3>
        <p className="mt-2 text-muted-foreground">
          Join the active ladder to unlock trading features and track your performance.
        </p>
      </Card>
    );
  }

  if (quote?.is_closed) {
    return (
      <Card
        data-testid="market-closed-guard"
        className="flex h-full min-h-[300px] flex-col items-center justify-center p-6 text-center"
      >
        <IconMoon className="mb-4 h-8 w-8 text-primary" />
        <h3 className="text-xl font-bold text-foreground">Market Closed</h3>
        <p className="mt-2 text-muted-foreground">
          Trading is currently unavailable.
          <br />
          Please come back during market hours.
        </p>
      </Card>
    );
  }

  if (!quote) {
    return (
      <Card
        data-testid="loading-market-guard"
        className="flex h-full min-h-[300px] flex-col items-center justify-center p-6 text-center"
      >
        <div className="mb-4 h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        <h3 className="text-xl font-bold text-foreground">Loading Market Data</h3>
        <p className="mt-2 text-muted-foreground">Fetching the latest quotes...</p>
      </Card>
    );
  }

  return <>{children}</>;
};
