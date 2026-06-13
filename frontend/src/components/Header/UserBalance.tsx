import { IconWallet } from '@/components/icons/CustomIcons';

interface UserBalanceProps {
  balance: number;
}

export const UserBalance = ({ balance }: UserBalanceProps) => (
  <div className="group flex items-center gap-2 text-muted-foreground bg-muted px-3 py-1.5 rounded-full border border-border">
    <IconWallet className="w-4 h-4 text-muted-foreground" />
    <span className="tabular-nums font-mono text-xs sm:text-sm">${balance.toFixed(2)}</span>
  </div>
);
