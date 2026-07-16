import { Link } from 'react-router-dom';
import styles from './Footer.module.css';

export const Footer = () => {
  return (
    <footer data-testid="app-footer" className={styles.footer}>
      <div className={styles.container}>
        <div className={styles.copyright}>
          &copy; {new Date().getFullYear()} Ticker Rush. All rights reserved.
        </div>
        <nav className={styles.links} aria-label="Footer navigation">
          <Link to="/impressum" className={styles.linkItem}>
            Impressum
          </Link>
          <Link to="/agb" className={styles.linkItem}>
            Terms (AGB)
          </Link>
          <Link to="/privacy" className={styles.linkItem}>
            Privacy Policy
          </Link>
          <a
            href="/api/swagger"
            target="_blank"
            rel="noopener noreferrer"
            className={styles.linkItem}
            aria-label="API Documentation (opens in a new tab)"
          >
            API Docs
          </a>
        </nav>
      </div>
    </footer>
  );
};
