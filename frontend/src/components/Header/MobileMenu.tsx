import { Navigation } from './Navigation';
import { Button } from '@/components/shared/Button';
import { IconLogOut } from '@icons/CustomIcons';
import { ThemeToggle } from '../ThemeToggle';

interface MobileMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onLogout: () => void;
  getLinkStyle: (isActive: boolean) => string;
}

export const MobileMenu = ({ isOpen, onClose, onLogout, getLinkStyle }: MobileMenuProps) => {
  if (!isOpen) return null;

  return (
    <div
      data-testid="mobile-menu"
      className="md:hidden fixed inset-0 z-40 bg-background/80 backdrop-blur-sm top-16"
    >
      <div className="bg-background border-b border-border p-4 shadow-lg space-y-4">
        <Navigation
          getLinkStyle={getLinkStyle}
          onItemClick={onClose}
          className="flex flex-col gap-2"
        />
        <div className="border-t border-border pt-4 flex items-center justify-between gap-4">
          <Button onClick={onLogout} variant="ghostDestructive" className="flex items-center gap-2">
            <IconLogOut className="w-4 h-4" />
            Logout
          </Button>
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-muted-foreground mr-2">Theme</span>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  );
};
