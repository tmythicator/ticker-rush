import { NavLink } from 'react-router-dom';
import { buttonVariants } from '@/components/shared/buttonVariants';

export const AuthButtons = () => (
  <div className="flex gap-2">
    <NavLink to="/login" className={buttonVariants({ variant: 'ghost' })}>
      Login
    </NavLink>
    <NavLink to="/register" className={buttonVariants({ variant: 'default' })}>
      Register
    </NavLink>
  </div>
);
