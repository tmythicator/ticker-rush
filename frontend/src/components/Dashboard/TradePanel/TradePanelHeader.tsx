import { IconTrending, IconRefresh } from '../../icons/CustomIcons';

export const TradePanelHeader = ({ isLoading }: { isLoading: boolean }) => {
  return (
    <div className="flex items-center justify-between mb-6">
      <h2 className="font-bold text-foreground flex items-center gap-2">
        <IconTrending className="w-4 h-4 text-primary" />
        Place Order
      </h2>
      {isLoading && <IconRefresh className="w-4 h-4 animate-spin text-muted-foreground" />}
    </div>
  );
};
