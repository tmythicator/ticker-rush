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
      <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center font-bold text-xs text-muted-foreground">
        {symbol[0] ?? '?'}
      </div>
      <div className="flex flex-col items-start text-foreground font-bold">
        <span className="text-sm font-bold">{symbol}</span>
        <div className="flex items-center gap-1.5 flex-wrap mt-0.5">
          <SourceBadge source={source} />
          {!isTradable && (
            <span
              data-testid="suspended-badge"
              className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-bold bg-destructive/10 text-destructive border border-destructive/20"
            >
              Suspended
            </span>
          )}
        </div>
      </div>
    </div>
  </td>
);
