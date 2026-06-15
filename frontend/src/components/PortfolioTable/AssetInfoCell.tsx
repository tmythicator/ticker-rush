import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerSource } from '@/types';

interface AssetInfoCellProps {
  symbol: string;
  source: TickerSource;
  isTradable?: boolean;
}

export const AssetInfoCell = ({ symbol, source, isTradable = true }: AssetInfoCellProps) => (
  <td className="px-6 py-4 font-bold text-foreground">
    <div className="flex items-center gap-3">
      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground">
        {symbol[0] ?? '?'}
      </div>
      <div className="flex flex-col items-start font-bold text-foreground">
        <span className="text-sm font-bold">{symbol}</span>
        <div className="mt-0.5 flex flex-wrap items-center gap-1.5">
          <SourceBadge source={source} />
          {!isTradable && (
            <span
              data-testid="suspended-badge"
              className="inline-flex items-center rounded border border-destructive/20 bg-destructive/10 px-1.5 py-0.5 text-[10px] font-bold text-destructive"
            >
              Suspended
            </span>
          )}
        </div>
      </div>
    </div>
  </td>
);
