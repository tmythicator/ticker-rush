import { SourceBadge } from '@/components/shared/SourceBadge';
import type { TickerSource } from '@/types';
import styles from './PortfolioTable.module.css';

interface AssetInfoCellProps {
  symbol: string;
  source: TickerSource;
  isTradable?: boolean;
}

export const AssetInfoCell = ({ symbol, source, isTradable = true }: AssetInfoCellProps) => (
  <td className={styles.cell}>
    <div className={styles.assetInfoWrapper}>
      <div className={styles.assetAvatar}>
        {symbol[0] ?? '?'}
      </div>
      <div className={styles.assetMeta}>
        <span className={styles.assetName}>{symbol}</span>
        <div className={styles.sourceBadges}>
          <SourceBadge source={source} />
          {!isTradable && (
            <span
              data-testid="suspended-badge"
              className={styles.suspendedBadge}
            >
              Suspended
            </span>
          )}
        </div>
      </div>
    </div>
  </td>
);
