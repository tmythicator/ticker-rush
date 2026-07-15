import { IconSettings } from '@icons/CustomIcons';
import { useState } from 'react';
import { Button } from '../shared/Button';
import { EditProfileModal } from './EditProfileModal';
import styles from './ProfileHeader.module.css';

export const ProfileHeader = () => {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  return (
    <header className={styles.header}>
      <div>
        <h1 className={styles.title}>My Profile</h1>
        <p className={styles.description}>Manage your assets and view performance</p>
      </div>
      <Button onClick={() => setIsEditModalOpen(true)} variant="secondary">
        <IconSettings className={styles.icon} />
        Edit Profile
      </Button>

      <EditProfileModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} />
    </header>
  );
};
