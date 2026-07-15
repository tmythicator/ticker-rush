import { NavLink } from 'react-router-dom';
import { IconBarChart, IconTrophy, IconUser } from '@/components/icons/CustomIcons';
import styles from './Header.module.css';

interface NavigationProps {
  onItemClick?: () => void;
  className?: string;
}

export const Navigation = ({ onItemClick, className }: NavigationProps) => (
  <nav className={className}>
    <NavLink
      to="/leaderboard"
      className={styles.navLink}
      onClick={onItemClick}
    >
      <IconTrophy />
      Ladder
    </NavLink>
    <NavLink
      to="/profile"
      className={styles.navLink}
      onClick={onItemClick}
    >
      <IconUser />
      Profile
    </NavLink>
    <NavLink
      to="/trade"
      className={styles.navLink}
      onClick={onItemClick}
    >
      <IconBarChart />
      Terminal
    </NavLink>
  </nav>
);
