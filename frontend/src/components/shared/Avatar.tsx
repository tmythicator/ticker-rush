import React from 'react';
import styles from './Avatar.module.css';

export interface AvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  initials: string;
  username?: string;
}

export const Avatar = React.forwardRef<HTMLDivElement, AvatarProps>(
  ({ initials, username, className, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={`${styles.avatar} ${className || ''}`}
        title={username}
        {...props}
      >
        {initials}
      </div>
    );
  },
);

Avatar.displayName = 'Avatar';
