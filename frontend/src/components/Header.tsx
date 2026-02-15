import {
  IconActivity,
  IconWallet,
  IconBarChart,
  IconTrophy,
  IconUser,
  IconLogOut,
} from '@icons/CustomIcons';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
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

  return (
    <header className="h-16 bg-background border-b border-border flex items-center px-4 lg:px-8 justify-between sticky top-0 z-50">
      <div className="flex items-center gap-8">
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

      <div className="flex items-center gap-4 text-sm font-medium">
        <ThemeToggle />
        {isAuthenticated && user ? (
          <>
            <div className="group flex items-center gap-2 text-muted-foreground bg-muted px-3 py-1.5 rounded-full border border-border">
              <IconWallet className="w-4 h-4 text-muted-foreground" />
              <span className="tabular-nums font-mono">${user.balance.toFixed(2)}</span>
            </div>
            <div
              className="w-9 h-9 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full border-2 border-background shadow-sm flex items-center justify-center text-white font-bold text-xs cursor-default"
              title={user.email}
            >
              {user.first_name ? user.first_name[0] : user.email[0]}
            </div>
            <button
              onClick={handleLogout}
              className="p-2 text-muted-foreground hover:text-foreground hover:bg-muted rounded-lg transition-colors"
              title="Logout"
            >
              <IconLogOut className="w-5 h-5" />
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
  );
};
