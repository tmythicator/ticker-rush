import { useQueryClient } from '@tanstack/react-query';
import { useCallback, type ReactNode } from 'react';
import { useUserQuery } from '../hooks/useUserQuery';
import { logout as apiLogout, type User } from '../lib/api';
import { QUERY_KEY_USER } from '../lib/queryKeys';
import { AuthContext } from './AuthContextDefinition';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient();
  const { data: user = null, isLoading } = useUserQuery();

  const setAuth = useCallback(
    (newUser: User | null) => {
      queryClient.setQueryData(QUERY_KEY_USER, newUser);
    },
    [queryClient],
  );

  const logout = useCallback(async () => {
    try {
      await apiLogout();
    } catch (error) {
      console.error('Logout failed', error);
    } finally {
      setAuth(null);
    }
  }, [setAuth]);

  return (
    <AuthContext.Provider
      value={{ isAuthenticated: !!user, isLoading: isLoading, logout, user, login: setAuth }}
    >
      {children}
    </AuthContext.Provider>
  );
};
