import { IconSettings } from '@icons/CustomIcons';
import { useState } from 'react';
import { Button } from '../shared/Button';
import { EditProfileModal } from './EditProfileModal';

export const ProfileHeader = () => {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  return (
    <header className="flex justify-between items-start">
      <div>
        <h1 className="text-3xl font-bold text-foreground">My Profile</h1>
        <p className="text-muted-foreground">Manage your assets and view performance</p>
      </div>
      <Button onClick={() => setIsEditModalOpen(true)} variant="secondary">
        <IconSettings className="w-4 h-4" />
        Edit Profile
      </Button>

      <EditProfileModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} />
    </header>
  );
};
