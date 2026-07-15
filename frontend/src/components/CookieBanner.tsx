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
    <div
      data-testid="cookie-banner"
      className={styles.overlay}
    >
      <div className={styles.card}>
        <div className={styles.textGroup}>
          <h2 className={styles.title}>Cookie Notice</h2>
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
    </div>
  );
};
