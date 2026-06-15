import React from 'react';
import { cn } from '@/lib/utils';

export interface AvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  initials: string;
  username?: string;
}

export const Avatar = React.forwardRef<HTMLDivElement, AvatarProps>(
  ({ initials, username, className, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(
          'flex h-9 w-9 cursor-pointer items-center justify-center rounded-full border-2 border-background bg-gradient-to-br from-blue-500 to-indigo-600 text-xs font-bold text-white shadow-sm transition-opacity hover:opacity-90',
          className,
        )}
        title={username}
        {...props}
      >
        {initials}
      </div>
    );
  },
);

Avatar.displayName = 'Avatar';
