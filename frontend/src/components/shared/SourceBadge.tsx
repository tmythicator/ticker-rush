import { type TickerSource } from '@/types';
import { cn } from '@/lib/utils';

export interface SourceBadgeProps extends React.ComponentProps<'span'> {
  source: TickerSource;
  ref?: React.Ref<HTMLSpanElement>;
}

export const SourceBadge = ({ source, className, ref, ...props }: SourceBadgeProps) => {
  const isCoinGecko = source === 'CoinGecko' || source === 'CG';
  const label = isCoinGecko ? 'Source: CoinGecko' : 'Source: Finnhub';
  const displayLabel = isCoinGecko ? 'CG' : 'FH';
  const colors = isCoinGecko ? 'bg-orange-500/20 text-orange-400' : 'bg-blue-500/20 text-blue-400';

  return (
    <span
      ref={ref}
      className={cn(
        'text-xs font-bold px-1.5 py-0.5 rounded cursor-help transition-colors',
        colors,
        className,
      )}
      title={label}
      {...props}
    >
      {displayLabel}
    </span>
  );
};
