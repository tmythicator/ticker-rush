import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/shared/Button';
import { Avatar } from '@/components/shared';
import { IconLogOut, IconMenu, IconX } from '@icons/CustomIcons';
import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { ThemeToggle } from './ThemeToggle/ThemeToggle';
import styles from './Header.module.css';

import { Logo } from './Logo';
import { Navigation } from './Navigation';
import { UserBalance } from './UserBalance';
import { AuthButtons } from './AuthButtons';
import { MobileMenu } from './MobileMenu';

export const Header = () => {
  const { user, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    setIsMobileMenuOpen(false);
    navigate('/login');
  };

  const toggleMenu = () => setIsMobileMenuOpen(!isMobileMenuOpen);

  return (
    <>
      <header data-testid="app-header" className={styles.header}>
        <div className={styles.navStart}>
          <Logo />
          {isAuthenticated && <Navigation className={styles.desktopNav} />}
        </div>

        <div className={styles.navEnd}>
          <div className={styles.smShow}>
            <ThemeToggle />
          </div>

          {isAuthenticated && user ? (
            <>
              <UserBalance balance={user.balance} />

              <NavLink to="/profile" className={styles.avatarLink}>
                <Avatar initials={user.first_name[0]} username={user.username} />
              </NavLink>
              <Button
                data-testid="logout-button"
                onClick={handleLogout}
                variant="ghost"
                size="icon"
                className={styles.desktopOnly}
                title="Logout"
              >
                <IconLogOut className={styles.logoutIcon} />
              </Button>
              <Button
                data-testid="mobile-menu-toggle"
                onClick={toggleMenu}
                variant="ghost"
                size="icon"
                className={styles.mdHide}
              >
                {isMobileMenuOpen ? (
                  <IconX className={styles.menuToggleIcon} />
                ) : (
                  <IconMenu className={styles.menuToggleIcon} />
                )}
              </Button>
            </>
          ) : (
            <AuthButtons />
          )}
        </div>
      </header>

      <MobileMenu
        isOpen={isMobileMenuOpen}
        onClose={() => setIsMobileMenuOpen(false)}
        onLogout={handleLogout}
      />
    </>
  );
};
