import { useAuth } from '@/hooks/useAuth';
import styles from './ProfilePage.module.css';
import { ProfileView } from '@/components/Profile/ProfileView';

export const ProfilePage = () => {
  const { user } = useAuth();

  if (!user) {
    return <div className={styles.loading}>Loading profile...</div>;
  }

  return <ProfileView user={user} />;
};
