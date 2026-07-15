import { IconArrowRight } from '@icons/CustomIcons';
import { Link, Navigate } from 'react-router-dom';
import { HomeChart } from './HomeChart';
import { useAuth } from '@/hooks/useAuth';
import buttonStyles from '@/components/shared/Button.module.css';
import styles from './HomePage.module.css';

export const HomePage = () => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return null;
  }

  if (user) {
    return <Navigate to="/profile" replace />;
  }

  return (
    <div className={styles.homeWrapper}>
      <section className={styles.heroSection}>
        <div className={styles.container}>
          <div className={styles.textGroup}>
            <h1 className={styles.title}>
              Play the Markets.
              <br />
              <span>Compete in Monthly Ladders.</span>
            </h1>
            <p className={styles.description}>
              Experience the thrill of real-time trading <strong>risk-free</strong>. Compete in{' '}
              <strong>monthly ladders</strong> and prove your skills to the world.
            </p>
          </div>
          <div className={styles.btnGroup}>
            <Link
              to="/register"
              className={`${buttonStyles.button} ${styles.btnItem}`}
              data-variant="default"
              data-size="lg"
            >
              Get Started <IconArrowRight className={styles.arrowIcon} />
            </Link>
            <Link
              to="/leaderboard"
              className={`${buttonStyles.button} ${styles.btnItem}`}
              data-variant="outline"
              data-size="lg"
            >
              View Leaderboard
            </Link>
          </div>
        </div>

        <div className={styles.chartWrapper}>
          <div className={styles.chartCard}>
            <HomeChart symbol="bitcoin" />
          </div>
        </div>
      </section>
    </div>
  );
};
