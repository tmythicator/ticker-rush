import { IconLock } from '@/components/icons/CustomIcons';

export const JoinLadderNotice = () => {
  return (
    <div className="mt-6 flex items-start gap-4 border-t border-border pt-6">
      <div className="mt-1">
        <IconLock className="h-4 w-4 text-secondary" />
      </div>
      <div className="space-y-2 text-xs">
        <p className="italic text-muted-foreground">
          <span className="font-bold text-secondary">Important:</span> Once you join, your
          participation in this ladder cycle is permanent and cannot be undone. This ensures the
          integrity of the leaderboard and fair competition.
        </p>
        <div className="flex items-center gap-2 text-primary/80">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <p>Privacy concern? You can always toggle your profile to Private in the settings.</p>
        </div>
      </div>
    </div>
  );
};
