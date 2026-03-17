import { IconMoon } from '@/components/icons/CustomIcons';
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
      <div className="bg-card rounded-lg shadow-sm border border-border p-6 flex flex-col items-center justify-center text-center h-full min-h-[300px]">
        <IconMoon className="w-8 h-8 mb-4 text-secondary/50" />
        <h3 className="text-xl font-bold text-foreground">Participation Required</h3>
        <p className="text-muted-foreground mt-2">
          Join the active ladder to unlock trading features and track your performance.
        </p>
      </div>
    );
  }

  if (quote?.is_closed) {
    return (
      <div className="bg-card rounded-lg shadow-sm border border-border p-6 flex flex-col items-center justify-center text-center h-full min-h-[300px]">
        <IconMoon className="w-8 h-8 mb-4 text-primary" />
        <h3 className="text-xl font-bold text-foreground">Market Closed</h3>
        <p className="text-muted-foreground mt-2">
          Trading is currently unavailable.
          <br />
          Please come back during market hours.
        </p>
      </div>
    );
  }

  if (!quote) {
    return (
      <div className="bg-card rounded-lg shadow-sm border border-border p-6 flex flex-col items-center justify-center text-center h-full min-h-[300px]">
        <div className="w-8 h-8 mb-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        <h3 className="text-xl font-bold text-foreground">Loading Market Data</h3>
        <p className="text-muted-foreground mt-2">Fetching the latest quotes...</p>
      </div>
    );
  }

  return <>{children}</>;
};
