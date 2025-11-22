import { type LucideIcon } from 'lucide-react';

interface StatCardProps {
    label: string;
    value: string;
    trend?: string;
    icon: LucideIcon;
}

export const StatCard = ({ label, value, trend, icon: Icon }: StatCardProps) => (
    <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex items-start justify-between hover:shadow-md transition-shadow">
        <div>
            <span className="text-xs text-slate-500 uppercase font-bold tracking-wider">{label}</span>
            <div className="text-2xl font-bold text-slate-900 mt-1">{value}</div>
            {trend && <div className="text-xs font-medium text-green-600 mt-1">{trend}</div>}
        </div>
        <div className="p-2 bg-slate-50 rounded-lg">
            <Icon className="w-4 h-4 text-slate-400" />
        </div>
    </div>
);