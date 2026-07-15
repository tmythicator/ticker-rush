import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerInfo, TickerSource } from '@/types';
import styles from './Leaderboard.module.css';

interface LeaderBoardAssetsProps {
  assets: TickerInfo[];
}

export const LeaderBoardAssets = ({ assets }: LeaderBoardAssetsProps) => {
  if (!assets || assets.length === 0) return null;

  return (
    <div className={styles.assetsSection}>
      <div className={styles.assetsHeader}>
        <span className={styles.assetsHeaderDot} />
        Tradable Assets
      </div>
      <div className={styles.assetsGrid}>
        {assets.map((t) => (
          <div
            key={t.symbol}
            className={styles.assetItem}
          >
            <span className={styles.assetSymbol}>{t.symbol}</span>
            <SourceBadge source={t.source as TickerSource} />
          </div>
        ))}
      </div>
    </div>
  );
};
