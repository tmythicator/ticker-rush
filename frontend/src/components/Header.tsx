import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/shared/Button';
import { buttonVariants } from '@/components/shared/buttonVariants';
import { Avatar } from '@/components/shared';
import { IconLogOut, IconMenu, IconX } from '@icons/CustomIcons';
import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { ThemeToggle } from './ThemeToggle';
import { cn } from '@/lib/utils';

import { Logo } from './Header/Logo';
import { Navigation } from './Header/Navigation';
import { UserBalance } from './Header/UserBalance';
import { AuthButtons } from './Header/AuthButtons';
import { MobileMenu } from './Header/MobileMenu';

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

  const getLinkStyle = (isActive: boolean): string => {
    return cn(
      buttonVariants({ variant: isActive ? 'default' : 'ghost' }),
      'flex items-center gap-2 border',
      !isActive && 'text-muted-foreground border-transparent hover:border-border',
    );
  };

  return (
    <>
      <header
        data-testid="app-header"
        className="sticky top-0 z-50 flex h-16 items-center justify-between border-b border-border bg-background px-4 lg:px-8"
      >
        <div className="flex items-center gap-4">
          <Logo />
          {isAuthenticated && (
            <Navigation getLinkStyle={getLinkStyle} className="hidden items-center gap-1 md:flex" />
          )}
        </div>

        <div className="flex items-center gap-3 text-sm font-medium md:gap-4">
          <div className="hidden sm:block">
            <ThemeToggle />
          </div>

          {isAuthenticated && user ? (
            <>
              <UserBalance balance={user.balance} />

              <NavLink to="/profile" className="block">
                <Avatar initials={user.first_name[0]} username={user.username} />
              </NavLink>
              <Button
                data-testid="logout-button"
                onClick={handleLogout}
                variant="ghost"
                size="icon"
                className="hidden md:flex"
                title="Logout"
              >
                <IconLogOut className="h-5 w-5" />
              </Button>
              <Button
                data-testid="mobile-menu-toggle"
                onClick={toggleMenu}
                variant="ghost"
                size="icon"
                className="md:hidden"
              >
                {isMobileMenuOpen ? (
                  <IconX className="h-5 w-5" />
                ) : (
                  <IconMenu className="h-5 w-5" />
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
        getLinkStyle={getLinkStyle}
      />
    </>
  );
};
