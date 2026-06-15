import { Link } from 'react-router-dom';

export const Footer = () => {
  return (
    <footer
      data-testid="app-footer"
      className="mt-auto w-full border-t border-border bg-background py-6"
    >
      <div className="container mx-auto flex flex-col items-center justify-between px-4 text-sm text-muted-foreground md:flex-row">
        <div className="mb-4 md:mb-0">
          &copy; {new Date().getFullYear()} Ticker Rush. All rights reserved.
        </div>
        <div className="flex space-x-6">
          <Link to="/impressum" className="transition-colors hover:text-primary">
            Impressum
          </Link>
          <Link to="/agb" className="transition-colors hover:text-primary">
            Terms (AGB)
          </Link>
          <Link to="/privacy" className="transition-colors hover:text-primary">
            Privacy Policy
          </Link>
        </div>
      </div>
    </footer>
  );
};
