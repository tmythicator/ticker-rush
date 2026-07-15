import { Link } from 'react-router-dom';
import { IconActivity } from '@/components/icons/CustomIcons';
import styles from './Header.module.css';

export const Logo = () => (
  <Link
    to="/"
    data-testid="header-logo"
    className={styles.logo}
  >
    <div className={styles.logoIcon}>
      <IconActivity />
    </div>
    <span className={styles.logoText}>
      Ticker Rush
    </span>
  </Link>
);
