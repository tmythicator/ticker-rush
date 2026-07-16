import { useState } from 'react';
import { Button } from '@/components/shared/Button';
import styles from './CookieBanner.module.css';

export const CookieBanner = () => {
  const [isVisible, setIsVisible] = useState(() => {
    return !localStorage.getItem('cookie-consent');
  });

  const handleAccept = () => {
    localStorage.setItem('cookie-consent', 'true');
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <aside data-testid="cookie-banner" className={styles.overlay} role="status" aria-live="polite">
      <div className={styles.card}>
        <div className={styles.textGroup}>
          <h3 className={styles.title}>Cookie Notice</h3>
          <p className={styles.description}>
            We use essential cookies for authentication and store game data (Redis/Postgres) to
            provide the service. You must agree to this usage to continue using the application.
          </p>
        </div>

        <div className={styles.footer}>
          <Button
            data-testid="cookie-banner-accept-button"
            onClick={handleAccept}
            className={styles.button}
          >
            I Understand & Agree
          </Button>
        </div>
      </div>
    </aside>
  );
};
