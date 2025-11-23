import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Header } from './components/Header';
import { Dashboard } from './pages/Dashboard';
import { Leaderboard } from './pages/Leaderboard';
import { Profile } from './pages/Profile';

function App() {
    return (
        <BrowserRouter>
            <div className="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans">
                <Header />

                <main className="flex-1 flex flex-col">
                    <Routes>
                        <Route path="/" element={<Leaderboard />} />
                        <Route path="/trade" element={<Dashboard />} />
                        <Route path="/profile" element={<Profile />} />
                        <Route path="*" element={<Navigate to="/" replace />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    );
}

export default App;