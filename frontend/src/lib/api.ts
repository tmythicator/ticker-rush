import { GetLeaderboardResponse, type Quote, type UpdateUserRequest, type User } from '@/types';

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
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const errorData = await res.json().catch(() => ({}));
      throw new ApiError(errorData.error || `Error: ${res.status}`, res.status);
    }
    return res.json();
  },
  put: async (endpoint: string, body: unknown) => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
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

export const buyStock = async (symbol: string, quantity: number): Promise<User> => {
  return api.post('/buy', { symbol, quantity });
};

export const sellStock = async (symbol: string, quantity: number): Promise<User> => {
  return api.post('/sell', { symbol, quantity });
};

export const login = async (username: string, password: string): Promise<User> => {
  return api.post('/login', { username, password });
};

export const logout = async (): Promise<void> => {
  return api.post('/logout', {});
};

export const getLeaderboard = async (limit = 10, offset = 0): Promise<GetLeaderboardResponse> => {
  const json = await api.get(`/leaderboard?limit=${limit}&offset=${offset}`);
  return GetLeaderboardResponse.fromJSON(json);
};

export const register = async (
  username: string,
  password: string,
  first_name: string,
  last_name: string,
): Promise<User> => {
  return api.post('/register', { username, password, first_name, last_name });
};

export const getConfig = async (): Promise<{ tickers: string[] }> => {
  return api.get('/config');
};

export const getHistory = async (symbol: string, limit = 100): Promise<Quote[]> => {
  return api.get(`/history?symbol=${symbol}&limit=${limit}`);
};

export const updateUser = async (data: UpdateUserRequest): Promise<User> => {
  return api.put('/user/me', data);
};
