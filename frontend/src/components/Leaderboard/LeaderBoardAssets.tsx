import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerInfo, TickerSource } from '@/types';

interface LeaderBoardAssetsProps {
  assets: TickerInfo[];
}

export const LeaderBoardAssets = ({ assets }: LeaderBoardAssetsProps) => {
  if (!assets || assets.length === 0) return null;

  return (
    <div className="mt-8 pt-6 border-t border-border/50">
      <div className="text-[10px] uppercase font-black text-muted-foreground tracking-widest mb-3 flex items-center gap-2">
        <span className="w-1.5 h-1.5 rounded-full bg-primary" />
        Tradable Assets
      </div>
      <div className="flex flex-wrap gap-2">
        {assets.map((t) => (
          <div
            key={t.symbol}
            className="flex items-center gap-2 bg-background px-3 py-1.5 rounded-lg border border-border hover:border-primary transition-colors group"
          >
            <span className="text-sm font-bold tracking-tight">{t.symbol}</span>
            <SourceBadge source={t.source as TickerSource} />
          </div>
        ))}
      </div>
    </div>
  );
};
