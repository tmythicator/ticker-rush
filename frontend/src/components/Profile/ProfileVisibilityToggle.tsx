import { Checkbox } from '@/components/shared/Checkbox';
import type { UseFormRegisterReturn } from 'react-hook-form';
import styles from './ProfileVisibilityToggle.module.css';
import { useId } from 'react';

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
  const inputId = useId();
  const descriptionId = useId();

  return (
    <label htmlFor={inputId} className={styles.container}>
      <div className={styles.textGroup}>
        <div className={styles.titleRow}>
          <span className={styles.title}>Profile Visibility</span>
          <span className={styles.badge}>{isPublic ? 'Public' : 'Private'}</span>
        </div>
        <p id={descriptionId} className={styles.description}>
          When public, your portfolio allocation is visible on the leaderboard.
        </p>
      </div>
      <Checkbox
        {...checkboxProps}
        id={inputId}
        aria-describedby={descriptionId}
        checked={isPublic}
        onChange={onToggle}
        data-testid="visibility-checkbox"
      />
    </label>
  );
};
