import { IconMoon, IconSun, IconSystem } from '@/components/icons/CustomIcons';
import { useIsMounted } from '@/hooks/useIsMounted';
import { useTheme } from 'next-themes';

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const mounted = useIsMounted();

  if (!mounted) {
    return (
      <div className="relative flex h-8 shrink-0 items-center gap-1 border-2 border-border bg-muted p-0.5 shadow-brutalist-sm">
        <div className="flex h-full w-7 items-center justify-center text-muted-foreground">
          <IconSun className="h-[14px] w-[14px]" />
        </div>
        <div className="flex h-full w-7 items-center justify-center text-muted-foreground">
          <IconSystem className="h-[14px] w-[14px]" />
        </div>
        <div className="flex h-full w-7 items-center justify-center text-muted-foreground">
          <IconMoon className="h-[14px] w-[14px]" />
        </div>
      </div>
    );
  }

  return (
    <div className="shadow-brutalist-sm relative flex h-8 shrink-0 items-center gap-1 border-2 border-border bg-muted p-0.5">
      {/* Light */}
      <button
        onClick={() => setTheme('light')}
        className={`relative z-10 flex h-full w-7 items-center justify-center transition-colors ${
          theme === 'light'
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:text-foreground'
        }`}
        aria-label="Light Mode"
      >
        <IconSun className="h-[14px] w-[14px]" />
      </button>

      {/* System (Default) */}
      <button
        onClick={() => setTheme('system')}
        className={`relative z-10 flex h-full w-7 items-center justify-center transition-colors ${
          theme === 'system'
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:text-foreground'
        }`}
        aria-label="System Mode"
      >
        <IconSystem className="h-[14px] w-[14px]" />
      </button>

      {/* Dark */}
      <button
        onClick={() => setTheme('dark')}
        className={`relative z-10 flex h-full w-7 items-center justify-center transition-colors ${
          theme === 'dark'
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:text-foreground'
        }`}
        aria-label="Dark Mode"
      >
        <IconMoon className="h-[14px] w-[14px]" />
      </button>
    </div>
  );
}
