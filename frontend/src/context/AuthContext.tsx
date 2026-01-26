import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { logout as apiLogout, getUser, type User } from '../lib/api';
import { AuthContext } from './AuthContextDefinition';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const login = useCallback((user: User) => {
    setUser(user);
  }, []);

  const logout = useCallback(async () => {
    try {
      await apiLogout();
    } catch (error) {
      console.error('Logout failed', error);
    } finally {
      setUser(null);
    }
  }, []);

  const refreshUser = useCallback((updatedUser: User) => {
    setUser(updatedUser);
  }, []);

  useEffect(() => {
    getUser()
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setIsLoading(false));
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, login, logout, refreshUser, isAuthenticated: !!user, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
};
