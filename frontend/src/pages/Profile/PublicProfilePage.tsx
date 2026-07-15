import { useParams } from 'react-router-dom';
import { PortfolioTable } from '@/components/PortfolioTable';
import { StatsGrid } from '@/components/Profile/StatsGrid';
import { IconLock, IconRefresh } from '@/components/icons/CustomIcons';
import { Card } from '@/components/shared/Card';
import { usePublicProfileQuery } from '@/hooks/usePublicProfileQuery';
import styles from './PublicProfilePage.module.css';

export const PublicProfilePage = () => {
  const { username } = useParams<{ username: string }>();

  const { data: user, isLoading, error } = usePublicProfileQuery(username);

  if (isLoading) {
    return (
      <div className={styles.loaderWrapper}>
        <IconRefresh className={styles.loaderIcon} />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div data-testid="profile-unavailable" className={styles.errorWrapper}>
        <IconLock className={styles.lockIcon} />
        <h1 className={styles.errorTitle}>Profile Unavailable</h1>
        <p className={styles.errorDescription}>This profile is private or does not exist.</p>
      </div>
    );
  }

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.glow1} />
      <div className={styles.glow2} />

      <div className={styles.container}>
        <div className={styles.headerGroup}>
          <h1 data-testid="profile-name" className={styles.title}>
            {user.first_name} {user.last_name}
          </h1>
          <div className={styles.metaGroup}>
            <span data-testid="profile-username" className={styles.username}>
              @{user.username}
            </span>
            {user.website && (
              <div className={styles.websiteGroup}>
                <span className={styles.websiteLabel}>Website:</span>
                <a
                  href={user.website}
                  target="_blank"
                  rel="noopener noreferrer"
                  className={styles.websiteLink}
                >
                  {user.website}
                </a>
              </div>
            )}
          </div>
        </div>

        <StatsGrid {...user} />

        <div className={styles.tableGroup}>
          <h2 className={styles.sectionTitle}>Portfolio</h2>
          <Card className={styles.tableContainer}>
            <PortfolioTable items={Object.values(user.portfolio || {})} isReadOnly />
          </Card>
        </div>
      </div>
    </div>
  );
};
