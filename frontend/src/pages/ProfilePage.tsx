import { PortfolioHoldings } from '@/components/PortfolioTable';
import { ProfileHeader, StatsGrid } from '@/components/Profile';
import { useAuth } from '@/hooks/useAuth';

export const ProfilePage = () => {
  const { user } = useAuth();

  if (!user) {
    return <div className="p-6">Loading profile...</div>;
  }

  return (
    <div className="mx-auto w-full max-w-6xl space-y-8 p-4 lg:p-6">
      <ProfileHeader />
      <StatsGrid {...user} />
      <PortfolioHoldings portfolio={user.portfolio ?? {}} />
    </div>
  );
};
