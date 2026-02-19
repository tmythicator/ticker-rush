import { PortfolioTable } from '@/components/PortfolioTable/PortfolioTable';
import { Header, StatsGrid } from '@/components/Profile';
import { useAuth } from '@/hooks/useAuth';

export const ProfilePage = () => {
  const { user } = useAuth();

  if (!user) {
    return <div className="p-6">Loading profile...</div>;
  }

  return (
    <div className="max-w-6xl w-full mx-auto p-4 lg:p-6 space-y-8">
      <Header />
      <StatsGrid {...user} />
      <PortfolioTable portfolio={user.portfolio ?? {}} />
    </div>
  );
};
