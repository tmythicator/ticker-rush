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
      <header className="h-16 bg-background border-b border-border flex items-center px-4 lg:px-8 justify-between sticky top-0 z-50">
        <div className="flex items-center gap-4">
          <Logo />
          {isAuthenticated && (
            <Navigation getLinkStyle={getLinkStyle} className="hidden md:flex items-center gap-1" />
          )}
        </div>

        <div className="flex items-center gap-3 md:gap-4 text-sm font-medium">
          <div className="hidden sm:block">
            <ThemeToggle />
          </div>

          {isAuthenticated && user ? (
            <>
              <UserBalance balance={user.balance} />

              <NavLink to="/profile" className="block">
                <Avatar initials={user.first_name[0]} username={user.username} />
              </NavLink>

              {/* Desktop Logout */}
              <Button
                onClick={handleLogout}
                variant="ghost"
                size="icon"
                className="hidden md:flex"
                title="Logout"
              >
                <IconLogOut className="w-5 h-5" />
              </Button>

              {/* Mobile Menu Toggle */}
              <Button onClick={toggleMenu} variant="ghost" size="icon" className="md:hidden">
                {isMobileMenuOpen ? (
                  <IconX className="w-5 h-5" />
                ) : (
                  <IconMenu className="w-5 h-5" />
                )}
              </Button>
            </>
          ) : (
            <AuthButtons />
          )}
        </div>
      </header>

      {/* Mobile Menu Overlay */}
      {isMobileMenuOpen && (
        <div className="md:hidden fixed inset-0 z-40 bg-background/80 backdrop-blur-sm top-16">
          <div className="bg-background border-b border-border p-4 shadow-lg space-y-4">
            <Navigation
              getLinkStyle={getLinkStyle}
              onItemClick={() => setIsMobileMenuOpen(false)}
              className="flex flex-col gap-2"
            />
            <div className="border-t border-border pt-4 flex items-center justify-between gap-4">
              <Button
                onClick={handleLogout}
                variant="ghostDestructive"
                className="flex items-center gap-2"
              >
                <IconLogOut className="w-4 h-4" />
                Logout
              </Button>
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium text-muted-foreground mr-2">Theme</span>
                <ThemeToggle />
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
