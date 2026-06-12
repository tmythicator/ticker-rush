import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerSource } from '@/types';

interface AssetInfoCellProps {
  symbol: string;
  source: TickerSource;
}

export const AssetInfoCell = ({ symbol, source }: AssetInfoCellProps) => (
  <td className="px-6 py-4 font-bold text-foreground">
    <div className="flex items-center gap-3">
      <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
        {symbol[0] ?? '?'}
      </div>
      <div className="flex flex-col items-start text-foreground font-bold">
        <span className="text-sm font-bold">{symbol}</span>
        <SourceBadge source={source} />
      </div>
    </div>
  </td>
);
