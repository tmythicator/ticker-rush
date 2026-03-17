import { IconTrophy } from '@/components/icons/CustomIcons';

interface LadderHeaderProps {
  name: string;
  type: string;
}

export const LadderHeader = ({ name, type }: LadderHeaderProps) => {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-primary font-bold uppercase tracking-widest text-xs bg-primary/10 w-fit px-3 py-1 rounded-full">
        <IconTrophy className="w-3.5 h-3.5" />
        Active Competition
      </div>
      <div className="space-y-2">
        <h1 className="text-4xl font-black tracking-tight">{name}</h1>
        <p className="text-muted-foreground max-w-xl leading-relaxed">
          Compete in this {type.toLowerCase()} ladder. The most profitable traders
          gain reputation and exclusive rewards!
        </p>
      </div>
    </div>
  );
};
