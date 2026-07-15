import { Checkbox } from '@/components/shared/Checkbox';
import type { UseFormRegisterReturn } from 'react-hook-form';
import styles from './ProfileVisibilityToggle.module.css';

interface ProfileVisibilityToggleProps {
  isPublic: boolean;
  onToggle: () => void;
  checkboxProps?: UseFormRegisterReturn;
}

export const ProfileVisibilityToggle = ({
  isPublic,
  onToggle,
  checkboxProps,
}: ProfileVisibilityToggleProps) => {
  return (
    <div className={styles.container} onClick={onToggle}>
      <div className={styles.textGroup}>
        <div className={styles.titleRow}>
          <span className={styles.title}>Profile Visibility</span>
          <span className={styles.badge}>{isPublic ? 'Public' : 'Private'}</span>
        </div>
        <p className={styles.description}>
          When public, your portfolio allocation is visible on the leaderboard.
        </p>
      </div>
      <Checkbox {...checkboxProps} data-testid="visibility-checkbox" />
    </div>
  );
};
