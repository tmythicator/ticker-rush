import { IconLock } from '@/components/icons/CustomIcons';

export const JoinLadderNotice = () => {
  return (
    <div className="mt-6 pt-6 border-t border-border flex items-start gap-4">
      <div className="mt-1">
        <IconLock className="w-4 h-4 text-secondary" />
      </div>
      <div className="text-xs space-y-2">
        <p className="text-muted-foreground italic">
          <span className="text-secondary font-bold">Important:</span> Once you join, your
          participation in this ladder cycle is permanent and cannot be undone. This ensures the
          integrity of the leaderboard and fair competition.
        </p>
        <div className="flex items-center gap-2 text-primary/80">
          <div className="w-1 h-1 bg-primary rounded-full" />
          <p>Privacy concern? You can always toggle your profile to Private in the settings.</p>
        </div>
      </div>
    </div>
  );
};
