import { PortfolioHoldings } from '@/components/PortfolioTable';
import { ProfileHeader, StatsGrid } from '@/components/Profile';
import { useAuth } from '@/hooks/useAuth';

export const ProfilePage = () => {
  const { user } = useAuth();

  if (!user) {
    return <div className="p-6">Loading profile...</div>;
  }

  return (
    <div className="max-w-6xl w-full mx-auto p-4 lg:p-6 space-y-8">
      <ProfileHeader />
      <StatsGrid {...user} />
      <PortfolioHoldings portfolio={user.portfolio ?? {}} />
    </div>
  );
};
