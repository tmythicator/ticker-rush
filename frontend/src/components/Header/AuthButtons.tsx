import { NavLink } from 'react-router-dom';
import buttonStyles from '@/components/shared/Button.module.css';
import styles from './Header.module.css';

export const AuthButtons = () => (
  <div className={styles.authButtons}>
    <NavLink
      to="/login"
      data-testid="login-link"
      className={buttonStyles.button}
      data-variant="ghost"
      data-size="sm"
    >
      Login
    </NavLink>
    <NavLink
      to="/register"
      data-testid="register-link"
      className={buttonStyles.button}
      data-variant="default"
      data-size="sm"
    >
      Register
    </NavLink>
  </div>
);
