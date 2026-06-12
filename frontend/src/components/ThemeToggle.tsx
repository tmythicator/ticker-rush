import { IconMoon, IconSun, IconSystem } from '@/components/icons/CustomIcons';
import { useIsMounted } from '@/hooks/useIsMounted';
import { useTheme } from 'next-themes';
import { ThemeToggleButton } from './ThemeToggleButton';

const THEME_OPTIONS = [
  { value: 'light', label: 'Light Mode', icon: IconSun },
  { value: 'system', label: 'System Mode', icon: IconSystem },
  { value: 'dark', label: 'Dark Mode', icon: IconMoon },
] as const;

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const mounted = useIsMounted();

  if (!mounted) {
    return (
      <div className="relative flex h-8 shrink-0 items-center gap-1 border-2 border-border bg-muted p-0.5 shadow-brutalist-sm">
        {THEME_OPTIONS.map(({ value, icon: Icon }) => (
          <div
            key={value}
            className="flex h-full w-7 items-center justify-center text-muted-foreground"
          >
            <Icon className="h-[14px] w-[14px]" />
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="shadow-brutalist-sm relative flex h-8 shrink-0 items-center gap-1 border-2 border-border bg-muted p-0.5">
      {THEME_OPTIONS.map(({ value, label, icon }) => (
        <ThemeToggleButton
          key={value}
          active={theme === value}
          onClick={() => setTheme(value)}
          icon={icon}
          label={label}
        />
      ))}
    </div>
  );
}
