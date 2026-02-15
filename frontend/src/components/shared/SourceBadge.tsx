import { type TickerSource } from '@/types';

interface SourceBadgeProps {
  source: TickerSource;
}

export const SourceBadge = ({ source }: SourceBadgeProps) => {
  const isCrypto = source === 'CG';
  const label = isCrypto ? 'Source: CoinGecko' : 'Source: Finnhub';
  const colors = isCrypto ? 'bg-orange-500/20 text-orange-400' : 'bg-blue-500/20 text-blue-400';

  return (
    <span className={`text-xs font-bold px-1.5 py-0.5 rounded cursor-help ${colors}`} title={label}>
      {source}
    </span>
  );
};
