import { Checkbox } from '@/components/shared/Checkbox';
import type { UseFormRegisterReturn } from 'react-hook-form';

interface ProfileVisibilityToggleProps {
  isPublic: boolean;
  onToggle: () => void;
  checkboxProps?: UseFormRegisterReturn;
}

const toggleCardStyles = {
  container:
    'bg-muted/30 p-4 rounded-xl border border-border/50 flex items-center justify-between group cursor-pointer hover:bg-muted/50 transition-colors',
  badge: 'text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary font-medium',
};

export const ProfileVisibilityToggle = ({
  isPublic,
  onToggle,
  checkboxProps,
}: ProfileVisibilityToggleProps) => {
  return (
    <div className={toggleCardStyles.container} onClick={onToggle}>
      <div className="space-y-1">
        <div className="flex items-center gap-2">
          <span className="font-medium text-foreground">Profile Visibility</span>
          <span className={toggleCardStyles.badge}>{isPublic ? 'Public' : 'Private'}</span>
        </div>
        <p className="text-xs text-muted-foreground">
          When public, your portfolio allocation is visible on the leaderboard.
        </p>
      </div>
      <Checkbox {...checkboxProps} data-testid="visibility-checkbox" />
    </div>
  );
};
