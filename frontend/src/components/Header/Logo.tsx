import { Link } from 'react-router-dom';
import { IconActivity } from '@/components/icons/CustomIcons';

export const Logo = () => (
  <Link to="/" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
    <div className="bg-primary p-1.5 rounded-lg shadow-sm">
      <IconActivity className="w-5 h-5 text-primary-foreground" />
    </div>
    <span className="font-bold text-lg tracking-tight text-foreground hidden sm:block">
      Ticker Rush
    </span>
  </Link>
);
