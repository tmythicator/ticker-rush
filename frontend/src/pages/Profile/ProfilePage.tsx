import { PortfolioHoldings } from '@/components/PortfolioTable';
import { ProfileHeader, StatsGrid } from '@/components/Profile';
import { useAuth } from '@/hooks/useAuth';
import styles from './ProfilePage.module.css';

export const ProfilePage = () => {
  const { user } = useAuth();

  if (!user) {
    return <div className={styles.loading}>Loading profile...</div>;
  }

  return (
    <div className={styles.profileWrapper}>
      <ProfileHeader />
      <StatsGrid {...user} />
      <PortfolioHoldings portfolio={user.portfolio ?? {}} />
    </div>
  );
};
