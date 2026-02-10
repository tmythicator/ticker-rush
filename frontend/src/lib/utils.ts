import { type PortfolioItem } from './api';
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Calculates the total invested capital from a user's portfolio.
 * Invested capital is the sum of (quantity * average_price) for all items.
 * Handles undefined/null portfolio gracefully.
 */
export const calculateInvestedCapital = (
  portfolio: Record<string, PortfolioItem> | undefined | null,
): number => {
  if (!portfolio) return 0;
  return Object.values(portfolio).reduce(
    (acc, item) => acc + item.quantity * item.average_price,
    0,
  );
};

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
