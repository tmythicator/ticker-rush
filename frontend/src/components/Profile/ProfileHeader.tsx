import { IconSettings } from '@icons/CustomIcons';
import { useId, useState } from 'react';
import { Button } from '../shared/Button';
import { EditProfileModal } from './EditProfileModal';
import styles from './ProfileHeader.module.css';
import type { PublicProfile, User } from '@/types';

interface ProfileHeaderProps {
  user: User | PublicProfile;
  isOwnProfile: boolean;
}

export const ProfileHeader = ({ user, isOwnProfile }: ProfileHeaderProps) => {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const descriptionId = useId();
  const profileDescription = isOwnProfile
    ? 'My Profile'
    : `${user.first_name || ''} ${user.last_name || ''}`.trim();

  return (
    <header className={styles.header}>
      <div className={styles.titleGroup}>
        <h1
          className={styles.title}
          data-testid="profile-name"
          aria-describedby={isOwnProfile ? descriptionId : undefined}
        >
          {profileDescription}
        </h1>
        {isOwnProfile ? (
          <p className={styles.description}>Manage your assets and view performance</p>
        ) : (
          <address className={styles.addressBlock}>
            <span className={styles.username} data-testid="profile-username">
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
                  aria-label={`Visit website: ${user.website} (opens in a new tab)`}
                >
                  {user.website}
                </a>
              </div>
            )}
          </address>
        )}
      </div>

      {isOwnProfile && (
        <>
          <Button
            onClick={() => setIsEditModalOpen(true)}
            variant="secondary"
            aria-haspopup="dialog"
            aria-expanded={isEditModalOpen}
          >
            <IconSettings className={styles.icon} aria-hidden="true" />
            Edit Profile
          </Button>

          <EditProfileModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} />
        </>
      )}
    </header>
  );
};
