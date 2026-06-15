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
            'rounded-xl p-2',
            isPnLPositive ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500',
          )}
        >
          <IconTrending className="h-5 w-5" />
        </div>
        <div>
          <span className="block text-xs font-bold uppercase tracking-wider text-muted-foreground">
            Total Gain/Loss
          </span>
          <div
            className={cn(
              'mt-0.5 text-2xl font-bold',
              isPnLPositive ? 'text-green-500' : 'text-red-500',
            )}
          >
            {formatCurrencyWithSign(totalPnL)}
          </div>
        </div>
      </div>
      <p className="mt-4 text-xs text-muted-foreground">
        Real-time P&L based on current market prices.
      </p>
    </Card>
  );
};
