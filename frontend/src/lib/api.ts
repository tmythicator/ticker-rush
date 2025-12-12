const API_URL = import.meta.env.VITE_API_URL;

let authToken: string | null = null;
export const setAuthToken = (token: string | null) => {
  authToken = token;
};

const getHeaders = (token?: string | null) => {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };
  const finalToken = token || authToken;
  if (finalToken) {
    headers['Authorization'] = `Bearer ${finalToken}`;
  }
  return headers;
};

export interface Quote {
  symbol: string;
  price: number;
  timestamp: number;
}

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export const api = {
  get: async (endpoint: string, token?: string | null) => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      headers: getHeaders(token),
    });
    if (!res.ok) throw new ApiError(`Error: ${res.status}`, res.status);
    return res.json();
  },
  post: async (endpoint: string, body: unknown, token?: string | null) => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method: 'POST',
      headers: getHeaders(token),
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const errorData = await res.json().catch(() => ({}));
      throw new ApiError(errorData.error || `Error: ${res.status}`, res.status);
    }
    return res.json();
  }
};

export const fetchQuote = async (symbol: string, token?: string | null): Promise<Quote> => {
  return api.get(`/quote?symbol=${symbol}`, token);
};

// TODO: use protobuf for frontend struct generation (autosync with backend)
export interface PortfolioItem {
  stock_symbol: string;
  quantity: number;
  average_price: number;
}

export interface User {
  email: string;
  first_name: string;
  last_name: string;
  balance: number;
  portfolio: Record<string, PortfolioItem>;
}

export const getUser = async (): Promise<User> => {
  return api.get(`/user/me`);
};

export const buyStock = async (symbol: string, count: number): Promise<User> => {
  return api.post('/buy', { symbol, count });
};

export const sellStock = async (symbol: string, count: number): Promise<User> => {
  return api.post('/sell', { symbol, count });
};

export const login = async (email: string, password: string): Promise<{ token: string; user: User }> => {
  return api.post('/login', { email, password });
};

export const register = async (email: string, password: string, first_name: string, last_name: string): Promise<User> => {
  return api.post('/register', { email, password, first_name, last_name });
};