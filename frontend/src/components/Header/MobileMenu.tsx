import { Navigation } from './Navigation';
import { Button } from '@/components/shared/Button';
import { IconLogOut } from '@icons/CustomIcons';
import { ThemeToggle } from './ThemeToggle/ThemeToggle';
import styles from './Header.module.css';

interface MobileMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onLogout: () => void;
  getLinkStyle: (isActive: boolean) => string;
}

export const MobileMenu = ({
  isOpen,
  onClose,
  onLogout,
  getLinkStyle,
}: MobileMenuProps) => {
  if (!isOpen) return null;

  return (
    <div
      data-testid="mobile-menu"
      className={styles.mobileOverlay}
    >
      <div className={styles.mobilePanel}>
        <Navigation
          getLinkStyle={getLinkStyle}
          onItemClick={onClose}
          className={styles.mobileNav}
        />
        <div className={styles.mobileFooter}>
          <Button
            onClick={onLogout}
            variant="ghostDestructive"
            className={styles.mobileLogoutBtn}
          >
            <IconLogOut />
            Logout
          </Button>
          <div className={styles.themeSelector}>
            <span className={styles.themeLabel}>Theme</span>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  );
};
