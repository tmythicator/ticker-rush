import { IconMoon, IconSun, IconSystem } from '@/components/icons/CustomIcons';
import { useIsMounted } from '@/hooks/useIsMounted';
import { useTheme } from 'next-themes';
import { ThemeToggleButton } from './ThemeToggleButton';
import styles from './ThemeToggle.module.css';

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
      <div className={styles.toggleContainer}>
        {THEME_OPTIONS.map(({ value, icon: Icon }) => (
          <div key={value} className={styles.togglePlaceholder}>
            <Icon />
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className={styles.toggleContainer}>
      {THEME_OPTIONS.map(({ value, label, icon }) => (
        <ThemeToggleButton
          key={value}
          active={theme === value}
          onClick={() => setTheme(value)}
          icon={icon}
          label={label}
          value={value}
        />
      ))}
    </div>
  );
}
