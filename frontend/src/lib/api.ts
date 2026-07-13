import {
  CreateTradeRequest,
  CreateTradeResponse,
  GetHistoryRequest,
  GetHistoryResponse,
  GetQuoteRequest,
  GetQuoteResponse,
  type Quote,
} from './proto/exchange/v1/exchange';
import { GetActiveLadderResponse, Ladder } from './proto/ladder/v1/ladder';
import { GetLeaderboardRequest, GetLeaderboardResponse } from './proto/leaderboard/v1/leaderboard';
import {
  CreateUserRequest,
  CreateUserResponse,
  GetMeResponse,
  GetPublicProfileRequest,
  GetPublicProfileResponse,
  LoginRequest,
  LoginResponse,
  UpdateUserRequest,
  UpdateUserResponse,
  type User,
  type PublicProfile,
} from './proto/user/v1/user';
import { ApiError, type InvalidParam, parseProblemDetails } from './errors';

const API_URL = `${import.meta.env.VITE_API_URL}/v1`;

export { ApiError, type InvalidParam };

const handleResponseError = async (res: Response): Promise<never> => {
  const contentType = res.headers.get('Content-Type') || '';
  if (
    contentType.includes('application/problem+json') ||
    contentType.includes('application/json')
  ) {
    const data = await res.json().catch(() => ({}));
    throw parseProblemDetails(data, res.status);
  }
  throw new ApiError(`Error: ${res.status}`, res.status);
};

export const api = {
  get: async <T>(endpoint: string): Promise<T> => {
    const res = await fetch(`${API_URL}${endpoint}`);
    if (!res.ok) {
      await handleResponseError(res);
    }
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
      await handleResponseError(res);
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
      await handleResponseError(res);
    }
    return res.json();
  },
  patch: async <T>(endpoint: string, body: unknown): Promise<T> => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      await handleResponseError(res);
    }
    return res.json();
  },
  delete: async <T = void>(endpoint: string): Promise<T> => {
    const res = await fetch(`${API_URL}${endpoint}`, {
      method: 'DELETE',
    });
    if (!res.ok) {
      await handleResponseError(res);
    }
    if (res.status === 204) return {} as T;
    return res.json().catch(() => ({}));
  },
};

export const getQuote = async (req: GetQuoteRequest): Promise<Quote> => {
  const json = await api.get(`/quotes/${req.symbol}`);
  const { quote } = GetQuoteResponse.fromJSON(json);
  if (!quote) throw new Error('Quote not found');
  return quote;
};

export const getUser = async (): Promise<User> => {
  const json = await api.get(`/profile`);
  const { user } = GetMeResponse.fromJSON(json);
  if (!user) throw new Error('User not found');
  return user;
};

export const createTrade = async (req: CreateTradeRequest): Promise<PublicProfile> => {
  const json = await api.post('/trades', req);
  const { participant } = CreateTradeResponse.fromJSON(json);
  if (!participant?.user) throw new Error('Update failed');
  return participant.user;
};

export const login = async (req: LoginRequest): Promise<User> => {
  const json = await api.post('/sessions', req);
  const { user } = LoginResponse.fromJSON(json);
  if (!user) throw new Error('Login failed');
  return user;
};

export const logout = async (): Promise<void> => {
  return api.delete('/sessions');
};

export const getLeaderboard = async (
  req: GetLeaderboardRequest,
): Promise<GetLeaderboardResponse> => {
  const json = await api.get(`/leaderboard?limit=${req.limit}&offset=${req.offset}`);
  return GetLeaderboardResponse.fromJSON(json);
};

export const register = async (req: CreateUserRequest): Promise<User> => {
  const json = await api.post('/users', req);
  const { user } = CreateUserResponse.fromJSON(json);
  if (!user) throw new Error('Registration failed');
  return user;
};

export const getActiveLadder = async (): Promise<Ladder | undefined> => {
  const json = await api.get('/ladder/active');
  const { ladder } = GetActiveLadderResponse.fromJSON(json);
  return ladder;
};

export const joinLadder = async (): Promise<void> => {
  await api.post('/ladder/participants', {});
};

export const getHistory = async (req: GetHistoryRequest): Promise<Quote[]> => {
  const json = await api.get(`/quotes/${req.symbol}/history?limit=${req.limit}`);
  const { history } = GetHistoryResponse.fromJSON(json);
  return history;
};

export const updateUser = async (req: UpdateUserRequest): Promise<User> => {
  const json = await api.patch('/profile', req);
  const { user } = UpdateUserResponse.fromJSON(json);
  if (!user) throw new Error('Update failed');
  return user;
};

export const getPublicProfile = async (req: GetPublicProfileRequest): Promise<PublicProfile> => {
  const json = await api.get(`/users/${req.username}`);
  const { profile } = GetPublicProfileResponse.fromJSON(json);
  if (!profile) throw new Error('Profile not found');
  return profile;
};
