import { useState, useEffect, useCallback, type ReactNode } from 'react';
import { setAuthToken, getUser, type User } from '../lib/api';
import { AuthContext } from './AuthContextDefinition';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user') || '') : null);
    const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
    const [isLoading] = useState(false);

    const login = useCallback((newToken: string, newUser: User) => {
        localStorage.setItem('token', newToken);
        localStorage.setItem('user', JSON.stringify(newUser));
        setToken(newToken);
        setUser(newUser);
        setAuthToken(newToken);
    }, []);

    const refreshUser = useCallback((updatedUser: User) => {
        localStorage.setItem('user', JSON.stringify(updatedUser));
        setUser(updatedUser);
    }, []);

    const logout = useCallback(() => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setToken(null);
        setUser(null);
        setAuthToken(null);
    }, []);

    useEffect(() => {
        if (token) {
            setAuthToken(token);
            getUser().then(refreshUser).catch(console.error);
        }
    }, [token, refreshUser]);


    return (
        <AuthContext.Provider value={{ user, token, login, logout, refreshUser, isAuthenticated: !!token, isLoading }}>
            {children}
        </AuthContext.Provider>
    );
};

