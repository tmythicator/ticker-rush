import { useState } from 'react';

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
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-background/80 backdrop-blur-sm p-4 animate-in fade-in duration-200">
      <div className="w-full max-w-lg bg-card border border-border rounded-xl shadow-2xl p-6 space-y-6">
        <div className="space-y-2">
          <h2 className="text-xl font-semibold tracking-tight">Cookie Notice</h2>
          <p className="text-muted-foreground">
            We use essential cookies for authentication and store game data (Redis/Postgres) to
            provide the service. You must agree to this usage to continue using the application.
          </p>
        </div>

        <div className="flex justify-end">
          <button
            onClick={handleAccept}
            className="bg-primary text-primary-foreground hover:bg-primary/90 px-8 py-3 rounded-md text-sm font-medium transition-colors w-full sm:w-auto"
          >
            I Understand & Agree
          </button>
        </div>
      </div>
    </div>
  );
};
