import { BaseAvatar } from '@/components/shared/BaseAvatar';
import type { User } from '@/types';

interface UserAvatarProps {
  user: User;
  className?: string;
}

export const UserAvatar = ({ user, className }: UserAvatarProps) => {
  const initials = user?.first_name?.[0]?.toUpperCase() || '?';
  const username = user?.username || 'Unknown';

  return <BaseAvatar initials={initials} label={username} className={className} />;
};
