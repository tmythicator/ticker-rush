import { IconWallet } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';

interface NetWorthCardProps {
  totalNetWorth: number;
  cash: number;
  assets: number;
}

export const NetWorthCard = ({ totalNetWorth, cash, assets }: NetWorthCardProps) => (
  <Card className="group relative overflow-hidden border-primary/20 bg-gradient-to-br from-primary/10 to-primary/5 p-6">
    <div className="absolute right-0 top-0 p-4 opacity-10 transition-opacity group-hover:opacity-20">
      <IconWallet className="h-24 w-24 text-primary" />
    </div>
    <span className="mb-2 block text-xs font-bold uppercase tracking-wider text-muted-foreground">
      Total Net Worth
    </span>
    <div className="font-mono text-4xl font-bold tracking-tight text-foreground">
      ${totalNetWorth.toFixed(2)}
    </div>
    <div className="mt-4 flex flex-wrap items-center gap-2 text-sm">
      <span className="rounded-full border border-border/50 bg-background/50 px-3 py-1.5 text-xs font-medium text-muted-foreground backdrop-blur-sm">
        Cash: <span className="font-semibold text-foreground">${cash.toFixed(2)}</span>
      </span>
      <span className="rounded-full border border-border/50 bg-background/50 px-3 py-1.5 text-xs font-medium text-muted-foreground backdrop-blur-sm">
        Assets: <span className="font-semibold text-foreground">${assets.toFixed(2)}</span>
      </span>
    </div>
  </Card>
);
