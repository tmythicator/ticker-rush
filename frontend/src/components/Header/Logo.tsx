import { Link } from 'react-router-dom';
import { IconActivity } from '@/components/icons/CustomIcons';
import styles from './Header.module.css';

export const Logo = () => (
  <Link
    to="/"
    data-testid="header-logo"
    className={styles.logo}
    aria-label="Ticker Rush Home"
  >
    <div className={styles.logoIcon} aria-hidden="true">
      <IconActivity />
    </div>
    <span className={styles.logoText}>Ticker Rush</span>
  </Link>
);
