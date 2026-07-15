import type { ComponentType, ComponentProps } from 'react';
import styles from './StatCard.module.css';

interface StatCardProps {
  label: string;
  value: string;
  trend?: string;
  icon: ComponentType<ComponentProps<'svg'>>;
}

export const StatCard = ({ label, value, trend, icon: Icon }: StatCardProps) => (
  <div className={styles.card}>
    <div>
      <span className={styles.label}>{label}</span>
      <div className={styles.value}>{value}</div>
      {trend && <div className={styles.trend}>{trend}</div>}
    </div>
    <div className={styles.iconWrapper}>
      <Icon className={styles.icon} />
    </div>
  </div>
);
