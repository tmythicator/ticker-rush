import { Link } from 'react-router-dom';

export const Footer = () => {
  return (
    <footer className="w-full py-6 mt-auto border-t border-border bg-background">
      <div className="container mx-auto px-4 flex flex-col md:flex-row justify-between items-center text-sm text-muted-foreground">
        <div className="mb-4 md:mb-0">
          &copy; {new Date().getFullYear()} Ticker Rush. All rights reserved.
        </div>
        <div className="flex space-x-6">
          <Link to="/impressum" className="hover:text-primary transition-colors">
            Impressum
          </Link>
          <Link to="/agb" className="hover:text-primary transition-colors">
            Terms (AGB)
          </Link>
          <Link to="/privacy" className="hover:text-primary transition-colors">
            Privacy Policy
          </Link>
        </div>
      </div>
    </footer>
  );
};
