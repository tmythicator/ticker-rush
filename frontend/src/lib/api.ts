import { GetConfigResponse } from './proto/config/v1/config';
import {
  BuyStockRequest,
  BuyStockResponse,
  GetHistoryRequest,
  GetHistoryResponse,
  GetQuoteRequest,
  GetQuoteResponse,
  SellStockRequest,
  SellStockResponse,
  type Quote,
} from './proto/exchange/v1/exchange';
import { GetLeaderboardRequest, GetLeaderboardResponse } from './proto/leaderboard/v1/leaderboard';
import {
  CreateUserRequest,
  CreateUserResponse,
  GetPublicProfileRequest,
  GetPublicProfileResponse,
  GetUserResponse,
  LoginRequest,
  LoginResponse,
  UpdateUserRequest,
  UpdateUserResponse,
  type User,
} from './proto/user/v1/user';

const API_URL = import.meta.env.VITE_API_URL;

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

export const api = {
  get: async <T>(endpoint: string): Promise<T> => {
    const res = await fetch(`${API_URL}${endpoint}`);
    if (!res.ok) throw new ApiError(`Error: ${res.status}`, res.status);
    return res.json();
  },
  post: async <T>(endpoint: string, body: unknown): Promise<T> => {
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
  put: async <T>(endpoint: string, body: unknown): Promise<T> => {
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

export const getQuote = async (req: GetQuoteRequest): Promise<Quote> => {
  const json = await api.get(`/quote?symbol=${req.symbol}`);
  const { quote } = GetQuoteResponse.fromJSON(json);
  if (!quote) throw new Error('Quote not found');
  return quote;
};

export const getUser = async (): Promise<User> => {
  const json = await api.get(`/user/me`);
  const { user } = GetUserResponse.fromJSON(json);
  if (!user) throw new Error('User not found');
  return user;
};

export const buyStock = async (req: BuyStockRequest): Promise<User> => {
  const json = await api.post('/buy', req);
  const { user } = BuyStockResponse.fromJSON(json);
  if (!user) throw new Error('Update failed');
  return user;
};

export const sellStock = async (req: SellStockRequest): Promise<User> => {
  const json = await api.post('/sell', req);
  const { user } = SellStockResponse.fromJSON(json);
  if (!user) throw new Error('Update failed');
  return user;
};

export const login = async (req: LoginRequest): Promise<User> => {
  const json = await api.post('/login', req);
  const { user } = LoginResponse.fromJSON(json);
  if (!user) throw new Error('Login failed');
  return user;
};

export const logout = async (): Promise<void> => {
  return api.post('/logout', {});
};

export const getLeaderboard = async (
  req: GetLeaderboardRequest,
): Promise<GetLeaderboardResponse> => {
  const json = await api.get(`/leaderboard?limit=${req.limit}&offset=${req.offset}`);
  return GetLeaderboardResponse.fromJSON(json);
};

export const register = async (req: CreateUserRequest): Promise<User> => {
  const json = await api.post('/register', req);
  const { user } = CreateUserResponse.fromJSON(json);
  if (!user) throw new Error('Registration failed');
  return user;
};

export const getConfig = async (): Promise<string[]> => {
  const json = await api.get('/config');
  const { tickers } = GetConfigResponse.fromJSON(json);
  return tickers;
};

export const getHistory = async (req: GetHistoryRequest): Promise<Quote[]> => {
  const json = await api.get(`/history?symbol=${req.symbol}&limit=${req.limit}`);
  const { history } = GetHistoryResponse.fromJSON(json);
  return history;
};

export const updateUser = async (req: UpdateUserRequest): Promise<User> => {
  const json = await api.put('/user/me', req);
  const { user } = UpdateUserResponse.fromJSON(json);
  if (!user) throw new Error('Update failed');
  return user;
};

export const getPublicProfile = async (req: GetPublicProfileRequest): Promise<User> => {
  const json = await api.get(`/users/${req.username}`);
  const { user } = GetPublicProfileResponse.fromJSON(json);
  if (!user) throw new Error('Profile not found');
  return user;
};
