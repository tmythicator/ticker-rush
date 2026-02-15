import { IconRefresh } from '@/components/icons/CustomIcons';

interface TradePanelHeaderProps {
  isLoading: boolean;
}

export const TradePanelHeader = ({ isLoading }: TradePanelHeaderProps) => {
  return (
    <div className="flex items-center justify-between mb-6">
      <div className="flex items-center gap-2">
        <h2 className="text-lg font-semibold text-foreground">Trade Asset</h2>
      </div>
      {isLoading && <IconRefresh className="w-4 h-4 animate-spin text-muted-foreground" />}
    </div>
  );
};
