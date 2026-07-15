import { IconRefresh } from '@/components/icons/CustomIcons';
import styles from './TradePanel.module.css';

interface TradePanelHeaderProps {
  isLoading: boolean;
}

export const TradePanelHeader = ({ isLoading }: TradePanelHeaderProps) => {
  return (
    <div className={styles.header}>
      <h2 className={styles.headerTitle}>Trade Asset</h2>
      {isLoading && <IconRefresh className={styles.spinner} />}
    </div>
  );
};
