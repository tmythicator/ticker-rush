import { cva } from 'class-variance-authority';

export const inputVariants = cva(
  'flex w-full rounded-lg border-2 border-input bg-background px-3 py-2 text-sm font-medium shadow-sm transition-all duration-100 placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:translate-x-[1px] focus-visible:translate-y-[1px] focus-visible:shadow-none disabled:cursor-not-allowed disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'border-input focus-visible:ring-ring',
        error: 'border-destructive focus-visible:ring-destructive',
        unstyled:
          'border-none bg-transparent p-0 shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 focus-visible:translate-x-0 focus-visible:translate-y-0 focus-visible:shadow-none h-auto w-auto',
      },
      size: {
        default: 'h-11',
        sm: 'h-9 px-2.5 text-xs',
        lg: 'h-12 px-4 text-base',
        unstyled: 'h-auto p-0',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
);
