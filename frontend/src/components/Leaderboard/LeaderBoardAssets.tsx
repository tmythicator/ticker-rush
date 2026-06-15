import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerInfo, TickerSource } from '@/types';

interface LeaderBoardAssetsProps {
  assets: TickerInfo[];
}

export const LeaderBoardAssets = ({ assets }: LeaderBoardAssetsProps) => {
  if (!assets || assets.length === 0) return null;

  return (
    <div className="mt-8 border-t border-border/50 pt-6">
      <div className="mb-3 flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-muted-foreground">
        <span className="h-1.5 w-1.5 rounded-full bg-primary" />
        Tradable Assets
      </div>
      <div className="flex flex-wrap gap-2">
        {assets.map((t) => (
          <div
            key={t.symbol}
            className="group flex items-center gap-2 rounded-lg border border-border bg-background px-3 py-1.5 transition-colors hover:border-primary"
          >
            <span className="text-sm font-bold tracking-tight">{t.symbol}</span>
            <SourceBadge source={t.source as TickerSource} />
          </div>
        ))}
      </div>
    </div>
  );
};
