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
          'w-9 h-9 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full border-2 border-background shadow-sm flex items-center justify-center text-white font-bold text-xs cursor-pointer hover:opacity-90 transition-opacity',
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
