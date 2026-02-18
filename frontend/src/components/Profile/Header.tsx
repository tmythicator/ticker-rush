import { useState } from 'react';
import { IconSettings } from '../icons/CustomIcons';
import { EditProfileModal } from './EditProfileModal';

export const Header = () => {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  return (
    <header className="flex justify-between items-start">
      <div>
        <h1 className="text-3xl font-bold text-foreground">My Profile</h1>
        <p className="text-muted-foreground">Manage your assets and view performance</p>
      </div>
      <button
        onClick={() => setIsEditModalOpen(true)}
        className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-foreground bg-muted hover:bg-muted/80 rounded-lg transition-colors border border-border"
      >
        <IconSettings className="w-4 h-4" />
        Edit Profile
      </button>

      <EditProfileModal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)} />
    </header>
  );
};
