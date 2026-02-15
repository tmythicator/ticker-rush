import { IconBriefcase, IconWallet, IconTrending } from '@icons/CustomIcons';

import { type User } from '@/types';
import { calculateInvestedCapital } from '@/lib/utils';

interface StatsGridProps {
  user: User;
}

export const StatsGrid = ({ user }: StatsGridProps) => {
  const portfolioItems = Object.values(user.portfolio ?? {});
  const investedCapital = calculateInvestedCapital(user.portfolio);
  const totalNetWorth = user.balance + investedCapital;
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div className="bg-primary text-primary-foreground p-6 rounded-lg shadow-sm relative overflow-hidden group">
        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <IconWallet className="w-24 h-24" />
        </div>
        <div className="text-primary-foreground/70 text-sm font-bold uppercase tracking-wider mb-2">
          Total Net Worth
        </div>
        <div className="text-4xl font-bold font-mono tracking-tight">
          ${totalNetWorth.toFixed(2)}
        </div>
        <div className="mt-4 flex flex-wrap items-center gap-2 text-sm">
          <span className="bg-primary-foreground/10 px-2 py-1 rounded-lg border border-primary-foreground/20">
            Cash: ${user.balance.toFixed(2)}
          </span>
          <span className="bg-primary-foreground/10 px-2 py-1 rounded-lg border border-primary-foreground/20">
            Assets: ${investedCapital.toFixed(2)}
          </span>
        </div>
      </div>

      <div className="bg-card p-6 rounded-lg shadow-sm border border-border">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-muted text-primary rounded-lg">
            <IconBriefcase className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground font-bold">Portfolio Items</div>
            <div className="text-2xl font-bold text-foreground">{portfolioItems.length}</div>
          </div>
        </div>
        <div className="text-sm text-muted-foreground">Active positions in your portfolio.</div>
      </div>

      <div className="bg-card p-6 rounded-lg shadow-sm border border-border">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 bg-muted text-green-500 rounded-lg">
            <IconTrending className="w-6 h-6" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground font-bold">Total Gain/Loss</div>
            <div className="text-2xl font-bold text-foreground">--</div>
          </div>
        </div>
        <div className="text-sm text-muted-foreground">
          Real-time P&L not available in this view.
        </div>
      </div>
    </div>
  );
};
