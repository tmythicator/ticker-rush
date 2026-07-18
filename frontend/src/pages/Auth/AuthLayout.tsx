import React from 'react';
import { Card } from '@/components/shared/Card';
import styles from './AuthLayout.module.css';

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
  footer: React.ReactNode;
}

export const AuthLayout = ({ title, subtitle, children, footer }: AuthLayoutProps) => {
  return (
    <div className={styles.pageWrapper}>
      <Card className={styles.card}>
        <div className={styles.header}>
          <h1 className={styles.title}>{title}</h1>
          {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
        </div>
        <div>{children}</div>
        <div className={styles.footer}>{footer}</div>
      </Card>
    </div>
  );
};
