import { IconTrending } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import { formatCurrencyWithSign } from '@/lib/utils';
import { cn } from '@/lib/utils';

interface TotalGainLossCardProps {
  totalPnL: number;
}

export const TotalGainLossCard = ({ totalPnL }: TotalGainLossCardProps) => {
  const isPnLPositive = totalPnL >= 0;

  return (
    <Card className="flex flex-col justify-between p-6">
      <div className="flex items-center gap-3">
        <div
          className={cn(
            'p-2 rounded-xl',
            isPnLPositive ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500',
          )}
        >
          <IconTrending className="w-5 h-5" />
        </div>
        <div>
          <span className="text-xs text-muted-foreground font-bold uppercase tracking-wider block">
            Total Gain/Loss
          </span>
          <div
            className={cn(
              'text-2xl font-bold mt-0.5',
              isPnLPositive ? 'text-green-500' : 'text-red-500',
            )}
          >
            {formatCurrencyWithSign(totalPnL)}
          </div>
        </div>
      </div>
      <p className="text-xs text-muted-foreground mt-4">
        Real-time P&L based on current market prices.
      </p>
    </Card>
  );
};
