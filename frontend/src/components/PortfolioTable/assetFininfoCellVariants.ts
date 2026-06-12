import { cva } from 'class-variance-authority';

export const assetFininfoCellVariants = cva('px-6 py-4 text-right font-mono', {
  variants: {
    variant: {
      default: 'text-foreground',
      muted: 'text-muted-foreground',
      medium: 'text-foreground font-medium',
      bold: 'text-foreground font-bold',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
});
