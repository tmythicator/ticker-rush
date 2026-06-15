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
      className="fixed inset-0 top-16 z-40 bg-background/80 backdrop-blur-sm md:hidden"
    >
      <div className="space-y-4 border-b border-border bg-background p-4 shadow-lg">
        <Navigation
          getLinkStyle={getLinkStyle}
          onItemClick={onClose}
          className="flex flex-col gap-2"
        />
        <div className="flex items-center justify-between gap-4 border-t border-border pt-4">
          <Button onClick={onLogout} variant="ghostDestructive" className="flex items-center gap-2">
            <IconLogOut className="h-4 w-4" />
            Logout
          </Button>
          <div className="flex items-center gap-2">
            <span className="mr-2 text-sm font-medium text-muted-foreground">Theme</span>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  );
};
