import { IconWallet } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';

interface NetWorthCardProps {
  totalNetWorth: number;
  cash: number;
  assets: number;
}

export const NetWorthCard = ({ totalNetWorth, cash, assets }: NetWorthCardProps) => (
  <Card className="relative overflow-hidden group bg-gradient-to-br from-primary/10 to-primary/5 border-primary/20 p-6">
    <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
      <IconWallet className="w-24 h-24 text-primary" />
    </div>
    <span className="text-muted-foreground text-xs font-bold uppercase tracking-wider block mb-2">
      Total Net Worth
    </span>
    <div className="text-4xl font-bold font-mono tracking-tight text-foreground">
      ${totalNetWorth.toFixed(2)}
    </div>
    <div className="mt-4 flex flex-wrap items-center gap-2 text-sm">
      <span className="bg-background/50 backdrop-blur-sm px-3 py-1.5 rounded-full border border-border/50 text-muted-foreground text-xs font-medium">
        Cash: <span className="text-foreground font-semibold">${cash.toFixed(2)}</span>
      </span>
      <span className="bg-background/50 backdrop-blur-sm px-3 py-1.5 rounded-full border border-border/50 text-muted-foreground text-xs font-medium">
        Assets: <span className="text-foreground font-semibold">${assets.toFixed(2)}</span>
      </span>
    </div>
  </Card>
);
