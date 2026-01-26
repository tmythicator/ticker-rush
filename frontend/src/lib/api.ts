import { Quote } from './proto/exchange/exchange';
import { PortfolioItem, User } from './proto/user/user';

export type { PortfolioItem, Quote, User };

const API_URL = import.meta.env.VITE_API_URL;

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export const api = {
  get: async (endpoint: string) => {
    const res = await fetch(`${API_URL}${endpoint}`);
    if (!res.ok) throw new ApiError(`Error: ${res.status}`, res.status);
    return res.json();
  },
  post: async (endpoint: string, body: unknown) => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method: 'POST',
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const errorData = await res.json().catch(() => ({}));
      throw new ApiError(errorData.error || `Error: ${res.status}`, res.status);
    }
    return res.json();
  },
};

export const fetchQuote = async (symbol: string): Promise<Quote> => {
  return api.get(`/quote?symbol=${symbol}`);
};

export const getUser = async (): Promise<User> => {
  return api.get(`/user/me`);
};

export const buyStock = async (symbol: string, count: number): Promise<User> => {
  return api.post('/buy', { symbol, count });
};

export const sellStock = async (symbol: string, count: number): Promise<User> => {
  return api.post('/sell', { symbol, count });
};

export const login = async (email: string, password: string): Promise<User> => {
  return api.post('/login', { email, password });
};

export const logout = async (): Promise<void> => {
  return api.post('/logout', {});
};

export const register = async (
  email: string,
  password: string,
  first_name: string,
  last_name: string,
): Promise<User> => {
  return api.post('/register', { email, password, first_name, last_name });
};
