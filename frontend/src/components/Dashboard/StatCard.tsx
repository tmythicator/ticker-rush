import type { ComponentType, ComponentProps } from 'react';

interface StatCardProps {
  label: string;
  value: string;
  trend?: string;
  icon: ComponentType<ComponentProps<'svg'>>;
}

export const StatCard = ({ label, value, trend, icon: Icon }: StatCardProps) => (
  <div className="flex items-start justify-between rounded-lg border border-border bg-card p-4 shadow-sm transition-shadow hover:shadow-md">
    <div>
      <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
        {label}
      </span>
      <div className="mt-1 text-2xl font-bold text-foreground">{value}</div>
      {trend && <div className="mt-1 text-xs font-medium text-green-600">{trend}</div>}
    </div>
    <div className="rounded-lg bg-muted p-2">
      <Icon className="h-4 w-4 text-muted-foreground" />
    </div>
  </div>
);
