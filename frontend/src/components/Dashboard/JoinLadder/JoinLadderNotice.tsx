import { IconLock } from '@/components/icons/CustomIcons';
import styles from './JoinLadderButton.module.css';

export const JoinLadderNotice = () => {
  return (
    <div className={styles.noticeContainer}>
      <div className={styles.noticeIconWrapper}>
        <IconLock className={styles.noticeIcon} />
      </div>
      <div className={styles.noticeTextGroup}>
        <p className={styles.noticeImportant}>
          <span className={styles.noticeImportantLabel}>Important:</span> Once you join, your
          participation in this ladder cycle is permanent and cannot be undone. This ensures the
          integrity of the leaderboard and fair competition.
        </p>
        <div className={styles.privacyBadge}>
          <div className={styles.privacyBadgeDot} />
          <p>Privacy concern? You can always toggle your profile to Private in the settings.</p>
        </div>
      </div>
    </div>
  );
};
