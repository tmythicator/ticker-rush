import * as React from 'react';
import styles from './Card.module.css';

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  ref?: React.Ref<HTMLDivElement>;
}

export const Card = ({ className, ref, ...props }: CardProps) => {
  return <div ref={ref} className={`${styles.card} ${className || ''}`} {...props} />;
};
