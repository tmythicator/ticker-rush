import { IconTrophy } from '@/components/icons/CustomIcons';

interface LadderHeaderProps {
  name: string;
  type: string;
}

export const LadderHeader = ({ name, type }: LadderHeaderProps) => {
  return (
    <div className="space-y-4">
      <div className="flex w-fit items-center gap-2 rounded-full bg-primary/10 px-3 py-1 text-xs font-bold uppercase tracking-widest text-primary">
        <IconTrophy className="h-3.5 w-3.5" />
        Active Competition
      </div>
      <div className="space-y-2">
        <h1 className="text-4xl font-black tracking-tight">{name}</h1>
        <p className="max-w-xl leading-relaxed text-muted-foreground">
          Compete in this {type.toLowerCase()} ladder. The most profitable traders gain reputation
          and exclusive rewards!
        </p>
      </div>
    </div>
  );
};
