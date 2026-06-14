import { NavLink } from 'react-router-dom';
import { buttonVariants } from '@/components/shared/buttonVariants';

export const AuthButtons = () => (
  <div className="flex gap-2">
    <NavLink to="/login" data-testid="login-link" className={buttonVariants({ variant: 'ghost' })}>
      Login
    </NavLink>
    <NavLink
      to="/register"
      data-testid="register-link"
      className={buttonVariants({ variant: 'default' })}
    >
      Register
    </NavLink>
  </div>
);
