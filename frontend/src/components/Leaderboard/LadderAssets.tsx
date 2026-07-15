import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerInfo, TickerSource } from '@/types';
import styles from './LadderAssets.module.css';
import { useId } from 'react';

interface LadderAssetsProps {
  assets: TickerInfo[];
}

export const LadderAssets = ({ assets }: LadderAssetsProps) => {
  if (!assets || assets.length === 0) return null;
  const headingId = useId();
  return (
    <section className={styles.assetsSection} aria-labelledby={headingId}>
      <div className={styles.assetsHeader}>
        <span className={styles.assetsHeaderDot} aria-hidden="true" />
        <h4 id={headingId} className={styles.assetsTitle}>
          Tradable Assets
        </h4>
      </div>

      <div className={styles.assetsGrid}>
        {assets.map((t) => (
          <div key={t.symbol} className={styles.assetItem}>
            <span className={styles.assetSymbol}>{t.symbol.toUpperCase()}</span>
            <SourceBadge source={t.source as TickerSource} className={styles.assetSourceBadge} />
          </div>
        ))}
      </div>
    </section>
  );
};
