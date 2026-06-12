import { IconBriefcase } from '@icons/CustomIcons';
import { Card } from '@/components/shared/Card';

interface PortfolioItemsCardProps {
  count: number;
}

export const PortfolioItemsCard = ({ count }: PortfolioItemsCardProps) => (
  <Card className="flex flex-col justify-between p-6">
    <div className="flex items-center gap-3">
      <div className="p-2 bg-primary/10 text-primary rounded-xl">
        <IconBriefcase className="w-5 h-5" />
      </div>
      <div>
        <span className="text-xs text-muted-foreground font-bold uppercase tracking-wider block">
          Portfolio Items
        </span>
        <div className="text-2xl font-bold text-foreground mt-0.5">{count}</div>
      </div>
    </div>
    <p className="text-xs text-muted-foreground mt-4">Active positions in your portfolio.</p>
  </Card>
);
