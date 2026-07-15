import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerInfo, TickerSource } from '@/types';
import styles from './LadderAssets.module.css';

interface LadderAssetsProps {
  assets: TickerInfo[];
}

export const LadderAssets = ({ assets }: LadderAssetsProps) => {
  if (!assets || assets.length === 0) return null;

  return (
    <div className={styles.assetsSection}>
      <div className={styles.assetsHeader}>
        <span className={styles.assetsHeaderDot} />
        Tradable Assets
      </div>
      <div className={styles.assetsGrid}>
        {assets.map((t) => (
          <div key={t.symbol} className={styles.assetItem}>
            <span className={styles.assetSymbol}>{t.symbol.toUpperCase()}</span>
            <SourceBadge source={t.source as TickerSource} className={styles.assetSourceBadge} />
          </div>
        ))}
      </div>
    </div>
  );
};
