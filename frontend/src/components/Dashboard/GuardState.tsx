import { Card } from '@/components/shared/Card';
import styles from './Guard.module.css';

interface GuardStateProps {
  icon?: React.ReactNode;
  title: string;
  description: React.ReactNode;
  testId?: string;
}

export const GuardState = ({ icon, title, description, testId }: GuardStateProps) => (
  <Card data-testid={testId} className={styles.card}>
    {icon}
    <h3 className={styles.title}>{title}</h3>
    <p className={styles.description}>{description}</p>
  </Card>
);
