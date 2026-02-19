import { useAuth } from '@/hooks/useAuth';
import {
  IconActivity,
  IconBarChart,
  IconLogOut,
  IconMenu,
  IconTrophy,
  IconUser,
  IconWallet,
  IconX,
} from '@icons/CustomIcons';
import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { ThemeToggle } from './ThemeToggle';

export const Header = () => {
  const { user, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getLinkStyle = (isActive: boolean): string => {
    const baseStyles =
      'flex items-center gap-2 px-4 py-2 text-sm font-bold transition-all border rounded-lg';

    if (isActive) {
      return `${baseStyles} bg-primary text-primary-foreground border-primary shadow-sm`;
    } else {
      return `${baseStyles} text-muted-foreground border-transparent hover:text-foreground hover:bg-muted hover:border-border`;
    }
  };

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const toggleMenu = () => setIsMobileMenuOpen(!isMobileMenuOpen);

  return (
    <>
      <header className="h-16 bg-background border-b border-border flex items-center px-4 lg:px-8 justify-between sticky top-0 z-50">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className="bg-primary p-1.5 rounded-lg shadow-sm">
              <IconActivity className="w-5 h-5 text-primary-foreground" />
            </div>
            <span className="font-bold text-lg tracking-tight text-foreground hidden sm:block">
              Ticker Rush
            </span>
          </div>

          {isAuthenticated && (
            <nav className="hidden md:flex items-center gap-1">
              <NavLink to="/leaderboard" className={(params) => getLinkStyle(params.isActive)}>
                <IconTrophy className="w-4 h-4" />
                Ladder
              </NavLink>
              <NavLink to="/profile" className={(params) => getLinkStyle(params.isActive)}>
                <IconUser className="w-4 h-4" />
                Profile
              </NavLink>
              <NavLink to="/trade" className={(params) => getLinkStyle(params.isActive)}>
                <IconBarChart className="w-4 h-4" />
                Terminal
              </NavLink>
            </nav>
          )}
        </div>

        <div className="flex items-center gap-3 md:gap-4 text-sm font-medium">
          <div className="hidden sm:block">
            <ThemeToggle />
          </div>

          {isAuthenticated && user ? (
            <>
              <div className="group flex items-center gap-2 text-muted-foreground bg-muted px-3 py-1.5 rounded-full border border-border">
                <IconWallet className="w-4 h-4 text-muted-foreground" />
                <span className="tabular-nums font-mono text-xs sm:text-sm">
                  ${user.balance.toFixed(2)}
                </span>
              </div>

              <NavLink to="/profile" className="block">
                <div
                  className="w-9 h-9 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full border-2 border-background shadow-sm flex items-center justify-center text-white font-bold text-xs cursor-pointer hover:opacity-90 transition-opacity"
                  title={user.username}
                >
                  {user.first_name[0]}
                </div>
              </NavLink>

              {/* Desktop Logout */}
              <button
                onClick={handleLogout}
                className="hidden md:block p-2 text-muted-foreground hover:text-foreground hover:bg-muted rounded-lg transition-colors"
                title="Logout"
              >
                <IconLogOut className="w-5 h-5" />
              </button>

              {/* Mobile Menu Toggle */}
              <button
                onClick={toggleMenu}
                className="md:hidden p-2 text-foreground hover:bg-muted rounded-lg transition-colors"
              >
                {isMobileMenuOpen ? (
                  <IconX className="w-5 h-5" />
                ) : (
                  <IconMenu className="w-5 h-5" />
                )}
              </button>
            </>
          ) : (
            <div className="flex gap-2">
              <NavLink
                to="/login"
                className="px-4 py-2 text-muted-foreground hover:text-foreground font-medium"
              >
                Login
              </NavLink>
              <NavLink
                to="/register"
                className="px-4 py-2 bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg font-medium transition-colors"
              >
                Register
              </NavLink>
            </div>
          )}
        </div>
      </header>

      {/* Mobile Menu Overlay */}
      {isMobileMenuOpen && (
        <div className="md:hidden fixed inset-0 z-40 bg-background/80 backdrop-blur-sm top-16">
          <div className="bg-background border-b border-border p-4 shadow-lg space-y-4">
            <nav className="flex flex-col gap-2">
              <NavLink
                to="/leaderboard"
                onClick={() => setIsMobileMenuOpen(false)}
                className={(params) => getLinkStyle(params.isActive)}
              >
                <IconTrophy className="w-4 h-4" />
                Ladder
              </NavLink>
              <NavLink
                to="/profile"
                onClick={() => setIsMobileMenuOpen(false)}
                className={(params) => getLinkStyle(params.isActive)}
              >
                <IconUser className="w-4 h-4" />
                Profile
              </NavLink>
              <NavLink
                to="/trade"
                onClick={() => setIsMobileMenuOpen(false)}
                className={(params) => getLinkStyle(params.isActive)}
              >
                <IconBarChart className="w-4 h-4" />
                Terminal
              </NavLink>
            </nav>
            <div className="border-t border-border pt-4 flex items-center justify-between gap-4">
              <button
                onClick={handleLogout}
                className="flex items-center gap-2 px-4 py-2 text-sm font-bold text-red-500 hover:bg-red-500/10 rounded-lg transition-colors"
              >
                <IconLogOut className="w-4 h-4" />
                Logout
              </button>
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
