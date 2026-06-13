import { NavLink } from 'react-router-dom';
import { IconBarChart, IconTrophy, IconUser } from '@/components/icons/CustomIcons';

interface NavigationProps {
  getLinkStyle: (isActive: boolean) => string;
  onItemClick?: () => void;
  className?: string;
}

export const Navigation = ({ getLinkStyle, onItemClick, className }: NavigationProps) => (
  <nav className={className}>
    <NavLink
      to="/leaderboard"
      className={(params) => getLinkStyle(params.isActive)}
      onClick={onItemClick}
    >
      <IconTrophy className="w-4 h-4" />
      Ladder
    </NavLink>
    <NavLink
      to="/profile"
      className={(params) => getLinkStyle(params.isActive)}
      onClick={onItemClick}
    >
      <IconUser className="w-4 h-4" />
      Profile
    </NavLink>
    <NavLink
      to="/trade"
      className={(params) => getLinkStyle(params.isActive)}
      onClick={onItemClick}
    >
      <IconBarChart className="w-4 h-4" />
      Terminal
    </NavLink>
  </nav>
);
