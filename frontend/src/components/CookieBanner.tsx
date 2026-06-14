import { useState } from 'react';
import { Button } from '@/components/shared/Button';

export const CookieBanner = () => {
  const [isVisible, setIsVisible] = useState(() => {
    return !localStorage.getItem('cookie-consent');
  });

  const handleAccept = () => {
    localStorage.setItem('cookie-consent', 'true');
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <div
      data-testid="cookie-banner"
      className="fixed inset-0 z-[100] flex items-center justify-center bg-background/80 backdrop-blur-sm p-4 animate-in fade-in duration-200"
    >
      <div className="w-full max-w-lg bg-card border border-border rounded-xl shadow-2xl p-6 space-y-6">
        <div className="space-y-2">
          <h2 className="text-xl font-semibold tracking-tight">Cookie Notice</h2>
          <p className="text-muted-foreground">
            We use essential cookies for authentication and store game data (Redis/Postgres) to
            provide the service. You must agree to this usage to continue using the application.
          </p>
        </div>

        <div className="flex justify-end">
          <Button
            data-testid="cookie-banner-accept-button"
            onClick={handleAccept}
            className="w-full sm:w-auto px-8"
          >
            I Understand & Agree
          </Button>
        </div>
      </div>
    </div>
  );
};
