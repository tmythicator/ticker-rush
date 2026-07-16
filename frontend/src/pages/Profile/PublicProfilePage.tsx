import { useParams } from 'react-router-dom';
import { IconLock, IconRefresh } from '@/components/icons/CustomIcons';
import { usePublicProfileQuery } from '@/hooks/usePublicProfileQuery';
import styles from './PublicProfilePage.module.css';
import { ProfileView } from '@/components/Profile/ProfileView';

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

  return <ProfileView user={user} />;
};
