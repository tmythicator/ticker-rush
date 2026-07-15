import { IconTrophy } from '@/components/icons/CustomIcons';
import styles from './LadderHeader.module.css';

interface LadderHeaderProps {
  name: string;
  type: string;
}

export const LadderHeader = ({ name, type }: LadderHeaderProps) => {
  return (
    <div className={styles.headerWrapper}>
      <div className={styles.badge}>
        <IconTrophy />
        Active Competition
      </div>
      <div className={styles.textGroup}>
        <h1 className={styles.title}>{name}</h1>
        <p className={styles.description}>
          Compete in this {type.toLowerCase()} ladder. The most profitable traders gain reputation
          and exclusive rewards!
        </p>
      </div>
    </div>
  );
};
