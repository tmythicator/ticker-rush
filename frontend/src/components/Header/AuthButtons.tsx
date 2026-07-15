import { NavLink } from 'react-router-dom';
import buttonStyles from '@/components/shared/Button.module.css';
import styles from './Header.module.css';
import clsx from 'clsx';

export const AuthButtons = () => (
  <div className={styles.authButtons}>
    <NavLink
      to="/login"
      data-testid="login-link"
      className={clsx(buttonStyles.button, buttonStyles.variantGhost, buttonStyles.sizeSm)}
    >
      Login
    </NavLink>
    <NavLink
      to="/register"
      data-testid="register-link"
      className={clsx(buttonStyles.button, buttonStyles.variantDefault, buttonStyles.sizeSm)}
    >
      Register
    </NavLink>
  </div>
);
