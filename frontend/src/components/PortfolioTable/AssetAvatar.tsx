import { BaseAvatar } from '@/components/shared/BaseAvatar';

interface AssetAvatarProps {
  symbol: string;
  className?: string;
}

export const AssetAvatar = ({ symbol, className }: AssetAvatarProps) => {
  const initials = symbol?.[0]?.toUpperCase() ?? '?';

  return (
    <BaseAvatar
      initials={initials}
      className={className}
      label={`${symbol} token`}
      aria-label={`Asset avatar for ${symbol}`}
    />
  );
};
