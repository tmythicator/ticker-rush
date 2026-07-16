import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/shared/Button';
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
import { UserAvatar } from './UserAvatar';

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
        <nav className={styles.navStart} aria-label="Main Navigation">
          <Logo />
          {isAuthenticated && <Navigation className={styles.desktopNav} />}
        </nav>

        <div className={styles.navEnd}>
          <div className={styles.smShow}>
            <ThemeToggle />
          </div>

          {isAuthenticated && user ? (
            <>
              <div className={styles.userMenu} role="group" aria-label="User menu">
                <UserBalance balance={user.balance} />
                <NavLink to="/profile" className={styles.avatarLink} aria-label="Go to profile">
                  <UserAvatar user={user} />
                </NavLink>
              </div>

              <Button
                data-testid="logout-button"
                onClick={handleLogout}
                variant="ghost"
                size="icon"
                className={styles.desktopOnly}
                aria-label="Logout"
              >
                <IconLogOut className={styles.logoutIcon} aria-hidden="true" />
              </Button>

              <Button
                data-testid="mobile-menu-toggle"
                onClick={toggleMenu}
                variant="ghost"
                size="icon"
                className={styles.mdHide}
                aria-expanded={isMobileMenuOpen}
                aria-label="Toggle mobile menu"
              >
                {isMobileMenuOpen ? (
                  <IconX className={styles.menuToggleIcon} aria-hidden="true" />
                ) : (
                  <IconMenu className={styles.menuToggleIcon} aria-hidden="true" />
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
