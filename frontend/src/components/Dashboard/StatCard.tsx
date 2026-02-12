import type { ComponentType, ComponentProps } from 'react';

interface StatCardProps {
  label: string;
  value: string;
  trend?: string;
  icon: ComponentType<ComponentProps<'svg'>>;
}

export const StatCard = ({ label, value, trend, icon: Icon }: StatCardProps) => (
  <div className="bg-card p-4 rounded-lg border border-border shadow-sm flex items-start justify-between hover:shadow-md transition-shadow">
    <div>
      <span className="text-xs text-muted-foreground uppercase font-bold tracking-wider">
        {label}
      </span>
      <div className="text-2xl font-bold text-foreground mt-1">{value}</div>
      {trend && <div className="text-xs font-medium text-green-600 mt-1">{trend}</div>}
    </div>
    <div className="p-2 bg-muted rounded-lg">
      <Icon className="w-4 h-4 text-muted-foreground" />
    </div>
  </div>
);
