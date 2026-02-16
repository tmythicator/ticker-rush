import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { type PortfolioItem, type TickerSource } from '@/types';

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

export function formatLocalTime(timestamp: number | string): string {
  if (!timestamp) return 'Never';
  return new Intl.DateTimeFormat('en-GB', {
    dateStyle: 'medium',
    timeStyle: 'medium',
  }).format(new Date(Number(timestamp) * 1000));
}

export function parseTicker(ticker: string): { source: TickerSource; symbol: string } {
  if (ticker.startsWith('FH:')) {
    return { source: 'FH', symbol: ticker.slice(3) };
  }
  if (ticker.startsWith('CG:')) {
    return { source: 'CG', symbol: ticker.slice(3) };
  }

  return { source: 'FH', symbol: ticker };
}

/**
 * Calculates the maximum quantity of an asset that can be purchased.
 *
 * @param buyingPower The user's available buying power.
 * @param price The current price of the asset.
 * @param percentage The fraction of buying power to use (default 1.0 = 100%).
 * @returns The quantity formatted to 6 decimal places.
 */
export const calculateMaxBuyQuantity = (
  buyingPower: number,
  price: number,
  percentage: number = 1.0,
): string => {
  if (!buyingPower || !price || price <= 0) return '0.000000';
  const adjustedPercentage = percentage >= 1.0 ? 0.999999 : percentage;
  return ((buyingPower * adjustedPercentage) / price).toFixed(6);
};
