import { IconWallet } from '@/components/icons/CustomIcons';

interface UserBalanceProps {
  balance: number;
}

export const UserBalance = ({ balance }: UserBalanceProps) => (
  <div className="group flex items-center gap-2 rounded-full border border-border bg-muted px-3 py-1.5 text-muted-foreground">
    <IconWallet className="h-4 w-4 text-muted-foreground" />
    <span className="font-mono text-xs tabular-nums sm:text-sm">${balance.toFixed(2)}</span>
  </div>
);
