import { IconBriefcase } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';

interface PortfolioItemsCardProps {
  count: number;
}

export const PortfolioItemsCard = ({ count }: PortfolioItemsCardProps) => (
  <Card className="flex flex-col justify-between p-6">
    <div className="flex items-center gap-3">
      <div className="rounded-xl bg-primary/10 p-2 text-primary">
        <IconBriefcase className="h-5 w-5" />
      </div>
      <div>
        <span className="block text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Portfolio Items
        </span>
        <div className="mt-0.5 text-2xl font-bold text-foreground">{count}</div>
      </div>
    </div>
    <p className="mt-4 text-xs text-muted-foreground">Active positions in your portfolio.</p>
  </Card>
);
