import { Link } from 'react-router-dom';
import { IconActivity } from '@/components/icons/CustomIcons';

export const Logo = () => (
  <Link
    to="/"
    data-testid="header-logo"
    className="flex items-center gap-2 transition-opacity hover:opacity-80"
  >
    <div className="rounded-lg bg-primary p-1.5 shadow-sm">
      <IconActivity className="h-5 w-5 text-primary-foreground" />
    </div>
    <span className="hidden text-lg font-bold tracking-tight text-foreground sm:block">
      Ticker Rush
    </span>
  </Link>
);
