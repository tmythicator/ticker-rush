import { createContext } from 'react';
import { type User } from '@/types';

export interface AuthContextType {
  user: User | null;
  login: (user: User | null) => void;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
