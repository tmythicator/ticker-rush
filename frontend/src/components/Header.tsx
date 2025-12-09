import { Activity, Wallet, BarChart2, Trophy, User, LogOut } from 'lucide-react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export const Header = () => {
    const { user, logout, isAuthenticated } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const getLinkStyle = (isActive: boolean): string => {
        const baseStyles = "flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors";

        if (isActive) {
            return `${baseStyles} bg-slate-100 text-blue-600`;
        } else {
            return `${baseStyles} text-slate-500 hover:text-slate-900 hover:bg-slate-50`;
        }
    };

    return (
        <header className="h-16 bg-white border-b border-slate-200 flex items-center px-4 lg:px-8 justify-between sticky top-0 z-50">
            <div className="flex items-center gap-8">
                <div className="flex items-center gap-2">
                    <div className="bg-blue-600 p-1.5 rounded-lg shadow-blue-100 shadow-sm">
                        <Activity className="w-5 h-5 text-white" />
                    </div>
                    <span className="font-bold text-lg tracking-tight text-slate-900 hidden sm:block">Ticker Rush</span>
                </div>

                {isAuthenticated && (
                    <nav className="hidden md:flex items-center gap-1">
                        <NavLink to="/" className={(params) => getLinkStyle(params.isActive)}>
                            <Trophy className="w-4 h-4" />
                            Ladder
                        </NavLink>
                        <NavLink to="/profile" className={(params) => getLinkStyle(params.isActive)}>
                            <User className="w-4 h-4" />
                            Profile
                        </NavLink>
                        <NavLink to="/trade" className={(params) => getLinkStyle(params.isActive)}>
                            <BarChart2 className="w-4 h-4" />
                            Terminal
                        </NavLink>
                    </nav>
                )}
            </div>

            <div className="flex items-center gap-4 text-sm font-medium">
                {isAuthenticated && user ? (
                    <>
                        <div className="group flex items-center gap-2 text-slate-600 bg-slate-50 px-3 py-1.5 rounded-full border border-slate-200">
                            <Wallet className="w-4 h-4 text-slate-400" />
                            <span className="tabular-nums font-mono">${user.balance.toFixed(2)}</span>
                        </div>
                        <div className="w-9 h-9 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full border-2 border-white shadow-sm flex items-center justify-center text-white font-bold text-xs cursor-default" title={user.email}>
                            {user.first_name ? user.first_name[0] : user.email[0]}
                        </div>
                        <button onClick={handleLogout} className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-50 rounded-lg transition-colors" title="Logout">
                            <LogOut className="w-5 h-5" />
                        </button>
                    </>
                ) : (
                    <div className="flex gap-2">
                        <NavLink to="/login" className="px-4 py-2 text-slate-600 hover:text-slate-900 font-medium">Login</NavLink>
                        <NavLink to="/register" className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors">Register</NavLink>
                    </div>
                )}
            </div>
        </header>
    );
};