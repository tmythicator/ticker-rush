import React from 'react';
import styles from './BaseAvatar.module.css';
import clsx from 'clsx';

export interface BaseAvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  initials: string;
  label?: string;
}

export const BaseAvatar = React.forwardRef<HTMLDivElement, BaseAvatarProps>(
  ({ initials, label, className, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={clsx(styles.avatar, className)}
        title={label}
        aria-label={label ? `Avatar of ${label}` : 'Avatar'}
        {...props}
      >
        {initials}
      </div>
    );
  },
);
