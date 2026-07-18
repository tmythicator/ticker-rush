import { PortfolioTable } from '@/components/PortfolioTable';
import { ProfileHeader, StatsGrid } from '@/components/Profile';
import { Card } from '@/components/shared';
import type { PublicProfile, User } from '@/types';
import styles from './ProfileView.module.css';
import { isOwnUser } from '@/lib/utils';
import { useId } from 'react';

interface ProfileViewProps {
  user: User | PublicProfile;
}

export const ProfileView = ({ user }: ProfileViewProps) => {
  const isOwnProfile = isOwnUser(user);
  const portfolioHeadingId = useId();
  const statsHeadingId = useId();

  return (
    <div className={styles.pageWrapper}>
      <div className={styles.ambientOrbPrimary} aria-hidden="true" />
      <div className={styles.ambientOrbSecondary} aria-hidden="true" />

      <article className={styles.container}>
        <ProfileHeader user={user} isOwnProfile={isOwnProfile} />

        <section aria-labelledby={statsHeadingId}>
          <h2 id={statsHeadingId} className="srOnly">
            {isOwnProfile ? 'My Financial Stats' : `${user.first_name}'s Financial Stats`}
          </h2>
          <StatsGrid {...user} />
        </section>

        <section className={styles.tableGroup} aria-labelledby={portfolioHeadingId}>
          <h2 id={portfolioHeadingId} className={styles.sectionTitle}>
            Current Holdings
          </h2>
          <Card className={styles.tableContainer}>
            <PortfolioTable
              items={Object.values(user.portfolio || {})}
              isReadOnly={!isOwnProfile}
            />
          </Card>
        </section>
      </article>
    </div>
  );
};
