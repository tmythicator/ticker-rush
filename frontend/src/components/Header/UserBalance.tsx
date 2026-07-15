import { IconWallet } from '@/components/icons/CustomIcons';
import styles from './Header.module.css';

interface UserBalanceProps {
  balance: number;
}

export const UserBalance = ({ balance }: UserBalanceProps) => (
  <div className={styles.balance}>
    <IconWallet />
    <span className={styles.balanceText}>${balance.toFixed(2)}</span>
  </div>
);
