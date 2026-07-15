import styles from './TradePanel.module.css';

export interface TradeFooterProps {
  buyingPower: number;
  estCost: number;
}

export const TradeFooter = ({ buyingPower, estCost }: TradeFooterProps) => {
  return (
    <div className={styles.footer}>
      <div className={styles.footerRow}>
        <span className={styles.footerLabel}>Buying Power</span>
        <span className={styles.footerValue}>${buyingPower.toFixed(2)}</span>
      </div>
      <div className={styles.footerRow}>
        <span className={styles.footerLabel}>Est. Cost</span>
        <span className={styles.footerEstCost}>${estCost.toFixed(2)}</span>
      </div>
    </div>
  );
};
