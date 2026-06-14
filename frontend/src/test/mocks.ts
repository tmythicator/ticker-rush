import type { User, Quote, PortfolioItem } from '@/types';

export const mockUserParticipating: User = {
  id: '1',
  username: 'testuser',
  first_name: 'Test',
  last_name: 'User',
  website: 'https://example.com',
  is_public: true,
  balance: 5000,
  portfolio: {
    AAPL: { stock_symbol: 'AAPL', quantity: 10, average_price: 150.0 },
  },
  is_participating: true,
  is_admin: false,
  is_banned: false,
  created_at: new Date('2026-06-01T00:00:00Z'),
};

export const mockUserNotParticipating: User = {
  ...mockUserParticipating,
  is_participating: false,
};

export const mockActiveQuote: Quote = {
  symbol: 'AAPL',
  price: 150.0,
  change: 1.5,
  change_percent: 1.0,
  timestamp: new Date().toISOString(),
  source: 'Finnhub',
  is_closed: false,
};

export const mockClosedQuote: Quote = {
  ...mockActiveQuote,
  is_closed: true,
};

export const mockPortfolioItemAAPL: PortfolioItem = {
  stock_symbol: 'AAPL',
  quantity: 10,
  average_price: 150.0,
};

export const mockPortfolioItemMSFT: PortfolioItem = {
  stock_symbol: 'MSFT',
  quantity: 5,
  average_price: 300.0,
};

export const mockPortfolio: Record<string, PortfolioItem> = {
  AAPL: mockPortfolioItemAAPL,
  MSFT: mockPortfolioItemMSFT,
};
